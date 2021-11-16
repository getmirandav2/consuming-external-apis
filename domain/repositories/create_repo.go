package repositories

import (
	"myapi/utils/errors"
	"strings"
)

type CreateRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *CreateRepoRequest) Validate() errors.APIError {
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return errors.NewBadRequestError("invalid repository name")
	}
	return nil
}

type CreateRepoResponse struct {
	ID    int    `json:"id"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type CreateReposResponse struct {
	StatusCode int                 `json:"status_code"`
	Results    []CreateReposResult `json:"results"`
}

type CreateReposResult struct {
	Response *CreateRepoResponse `json:"repository"`
	Error    errors.APIError     `json:"error"`
}
