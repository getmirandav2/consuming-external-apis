package github_provider

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"myapi/clients/restclient"
	"myapi/domain/github"
	"net/http"
)

const (
	headerAuthorization       = "Authorization"
	headerAuthorizationFormat = "token %s"

	urlCreateRepo = "https://api.github.com/user/repos"
)

func getHeaderAuthorization(token string) string {
	return fmt.Sprintf(headerAuthorizationFormat, token)
}

func CreateRepo(token string, request *github.CreateRepoRequest) (*github.CreateRepoResponse, *github.GitHubErrorResponse) {
	headers := &http.Header{}
	headers.Set(headerAuthorization, getHeaderAuthorization(token))
	response, err := restclient.Post(urlCreateRepo, request, headers)
	if err != nil {
		log.Println("error when trying o create new repo in github:", err.Error())
		return nil, &github.GitHubErrorResponse{StatusCode: http.StatusInternalServerError, Message: err.Error()}
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, &github.GitHubErrorResponse{StatusCode: http.StatusInternalServerError, Message: "invalid response body"}
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {
		var errResponse github.GitHubErrorResponse
		if err := json.Unmarshal(bytes, &errResponse); err != nil {
			return nil, &github.GitHubErrorResponse{StatusCode: http.StatusInternalServerError, Message: "invalid json error response body"}
		}
		errResponse.StatusCode = response.StatusCode
		return nil, &errResponse
	}

	var result github.CreateRepoResponse
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println("error when trying to unmarshal create repo successful response:", err.Error())
		return nil, &github.GitHubErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error when trying to unmarshal github create repo response"}
	}
	return &result, nil
}
