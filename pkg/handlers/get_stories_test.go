package handlers

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	"github.com/snigdhasambitak/hackernews-api/pkg/models"
)

func TestGetStories(t *testing.T) {
	// create mock instance of mockHN service
	mockService := &mockHN{}
	expected := []models.Story{
		{Author: "yl", Karma: 20, Comments: 10, Title: "First", Position: 3},
		{Author: "sa", Karma: 30, Comments: 18, Title: "Second", Position: 2},
		{Author: "abd", Karma: 25, Comments: 25, Title: "third", Position: 1},
		{Author: "xyz", Karma: 77, Comments: 9, Title: "fourth", Position: 4},
		{Author: "dd", Karma: 90, Comments: 7, Title: "fifth", Position: 5},
	}
	// mock func Curated50 with required arguments and expected return arguments
	mockService.On("Curated50", 2413).Return(expected, nil)
	mockRW := &mockResponseWriter{}
	// mock func WriteHeader with http.StatusOK argument
	mockRW.On("WriteHeader", http.StatusOK).Return()
	mockRW.On("Header").Return(make(http.Header, 0))
	mockRW.On("Write", mock.Anything).Return(1, nil)

	c := &gin.Context{Writer: mockRW}
	// initialize handlers with mockHN service
	h := &handlers{hn: mockService}
	h.GetStories(c)
	// assert all the mocks have been called
	mockService.AssertExpectations(t)
	mockRW.AssertExpectations(t)
}
