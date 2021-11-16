package restclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	enableMocks = false
	mocks       = make(map[string]*Mock)
)

type Mock struct {
	URL        string
	HttpMethod string
	Response   *http.Response
	Err        error
}

func StartMockups() {
	enableMocks = true
}

func StopMockups() {
	enableMocks = false
}

func getMockId(httpMethod, url string) string {
	return fmt.Sprintf("%s_%s", httpMethod, url)
}

func AddMockup(mock Mock) {
	mocks[getMockId(mock.HttpMethod, mock.URL)] = &mock
}

func FlushMockups() {
	mocks = make(map[string]*Mock)
}

func Post(url string, body interface{}, headers *http.Header) (*http.Response, error) {
	if enableMocks {
		mock := mocks[getMockId(http.MethodPost, url)]
		if mock == nil {
			return nil, errors.New("no mockup found for give url")
		}
		return mock.Response, mock.Err
	}
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	request.Header = *headers

	client := http.Client{}
	return client.Do(request)
}
