package main

import (
	"bufio"
	"fmt"
	"log"
	"myapi/domain/repositories"
	"myapi/services"
	"myapi/utils/errors"
	"os"
	"sync"
)

var (
	success map[string]string
	failure map[string]errors.APIError
)

type createRepoResult struct {
	Request *repositories.CreateRepoRequest
	Result  *repositories.CreateRepoResponse
	Error   errors.APIError
}

func getRequests() []repositories.CreateRepoRequest {
	result := make([]repositories.CreateRepoRequest, 0)

	file, err := os.Open("requests.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		request := repositories.CreateRepoRequest{
			Name: scanner.Text(),
		}
		result = append(result, request)
	}
	fmt.Println(result)

	return result
}

func main() {
	requests := getRequests()

	fmt.Printf("%d requests\n", len(requests))

	input := make(chan createRepoResult)
	buffer := make(chan bool, 10)
	var wg sync.WaitGroup

	go handleResults(&wg, input)

	for _, request := range requests {
		fmt.Printf("%s\n", request.Name)
		buffer <- true
		wg.Add(1)
		go createRepo(request, input, buffer)
	}

	wg.Wait()
	close(input)
}

func handleResults(wg *sync.WaitGroup, input chan createRepoResult) {
	success = make(map[string]string)
	failure = make(map[string]errors.APIError)

	for result := range input {
		if result.Error != nil {
			failure[result.Request.Name] = result.Error
		} else {
			success[result.Request.Name] = result.Result.Name
		}
		wg.Done()
	}
}

func createRepo(request repositories.CreateRepoRequest, output chan createRepoResult, buffer chan bool) {
	result, err := services.RepositoryService.CreateRepo(&request)
	output <- createRepoResult{
		Request: &request,
		Result:  result,
		Error:   err,
	}
	<-buffer
}
