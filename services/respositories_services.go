package services

import (
	"myapi/config"
	"myapi/domain/github"
	"myapi/domain/repositories"
	"myapi/providers/github_provider"
	"myapi/utils/errors"
	"net/http"
	"sync"
)

type reposService struct{}

type reposServiceInterface interface {
	CreateRepo(input *repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.APIError)
	CreateRepos(input []repositories.CreateRepoRequest) (*repositories.CreateReposResponse, errors.APIError)
}

var (
	RepositoryService reposServiceInterface
)

func init() {
	RepositoryService = &reposService{}
}

func (r *reposService) CreateRepo(input *repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.APIError) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	request := github.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Private:     true,
	}

	response, err := github_provider.CreateRepo(config.GetGithubAccessToken(), &request)
	if err != nil {
		return nil, errors.NewAPIError(err.StatusCode, err.Message)
	}
	// arma la respuesta fuera de la logica de como crear un repo en github
	result := repositories.CreateRepoResponse{
		ID:    response.ID,
		Name:  response.Name,
		Owner: response.Owner.Login,
	}
	return &result, nil
}

func (s *reposService) CreateRepos(requests []repositories.CreateRepoRequest) (*repositories.CreateReposResponse, errors.APIError) {
	input := make(chan repositories.CreateReposResult)
	output := make(chan repositories.CreateReposResponse)
	defer close(output)

	var wg sync.WaitGroup
	go s.handleRepoResults(&wg, input, output)

	for _, current := range requests {
		wg.Add(1)
		go s.createRepoConcurrent(current, input)
	}

	wg.Wait()
	close(input)

	result := <-output

	successCreations := 0
	for _, result := range result.Results {
		if result.Error == nil {
			successCreations++
		}
	}

	if successCreations == 0 {
		result.StatusCode = result.Results[0].Error.Status()
	} else if successCreations == len(requests) {
		result.StatusCode = http.StatusCreated
	} else {
		result.StatusCode = http.StatusPartialContent
	}
	return &result, nil
}

func (r *reposService) handleRepoResults(wg *sync.WaitGroup, input chan repositories.CreateReposResult, ouput chan repositories.CreateReposResponse) {
	var results repositories.CreateReposResponse
	// Hasta que cerremos el canal saldrá de esta for.
	for incommingEvent := range input {
		repoResult := &repositories.CreateReposResult{
			Response: incommingEvent.Response,
			Error:    incommingEvent.Error,
		}
		results.Results = append(results.Results, *repoResult)
		wg.Done()
	}
	// Cuando el canal esté cerrado, esto se ejecuta.
	ouput <- results
}

func (s *reposService) createRepoConcurrent(input repositories.CreateRepoRequest, ouput chan repositories.CreateReposResult) {
	if err := input.Validate(); err != nil {
		ouput <- repositories.CreateReposResult{Error: err}
		return
	}
	result, err := s.CreateRepo(&input)
	if err != nil {
		ouput <- repositories.CreateReposResult{Error: err}
		return
	}
	ouput <- repositories.CreateReposResult{Response: result}
}
