package test_utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetMockedContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	request, err := http.NewRequest(http.MethodGet, "http://localhost:123/something", nil)
	assert.Nil(t, err)
	response := httptest.NewRecorder()
	request.Header = http.Header{"X-Mock": {"true"}}
	c := GetMockedContext(request, response)

	assert.EqualValues(t, "true", c.GetHeader("X-Mock"))
	assert.EqualValues(t, http.MethodGet, c.Request.Method)
	assert.EqualValues(t, "/something", c.Request.URL.Path)
}
