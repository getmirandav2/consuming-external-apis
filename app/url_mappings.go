package app

import "myapi/controllers/repositories"

func MapUrls() {
	router.POST("/repository", repositories.CreateRepo)
	router.POST("/repositories", repositories.CreateRepos)
}
