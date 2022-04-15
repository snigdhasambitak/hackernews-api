package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	"github.com/snigdhasambitak/hackernews-api/pkg/hackernews"
	"github.com/snigdhasambitak/hackernews-api/pkg/models"
)

// mockHN is mock implementation of hackernews.Service
type mockHN struct {
	hackernews.Service
	mock.Mock
}

// mockResponseWriter is mock implementation of http.ResponseWriter
type mockResponseWriter struct {
	gin.ResponseWriter
	mock.Mock
}

// Curated50 mocks the call and returns predefined arguments
func (m *mockHN) Curated50(minKarma int) ([]models.Story, error) {
	// register the call to method which returns the return arguments for mock
	args := m.Called(minKarma)
	// cast return arguments to match the functions expected return arguments
	return args.Get(0).([]models.Story), args.Error(1)
}

func (m *mockResponseWriter) JSON(code int, obj interface{}) {
	_ = m.Called(code, obj)
	return
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	_ = m.Called(statusCode)
	return
}

func (m *mockResponseWriter) Header() http.Header {
	args := m.Called()
	return args.Get(0).(http.Header)
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	args := m.Called(data)
	return args.Int(0), args.Error(1)
}
