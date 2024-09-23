package mock

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	mocking "github.com/stretchr/testify/mock"
)

func TestRouterRun(t *testing.T) {
	srv := serviceMock{}
	router := newRouter(&srv)
	assert.NotNil(t, router)
	assert.NotPanics(t, func() {
		go router.Run(":8080")
	})
}

func TestMappingEntrypointErrorMethodNotAllowed(t *testing.T) {
	srv := serviceMock{}
	router := newRouter(&srv)
	response := responseWriterMock{}
	response.On("WriteHeader", http.StatusMethodNotAllowed).Return(nil)
	request := http.Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost",
			Path:   "/mock/mapping",
		},
		Method: http.MethodGet,
	}
	router.server.ServeHTTP(&response, &request)
	response.AssertCalled(t, "WriteHeader", http.StatusMethodNotAllowed)
	response.AssertExpectations(t)
}

func TestMappingEntrypointErrorDecodeJson(t *testing.T) {
	body := io.NopCloser(strings.NewReader("{invalid json}"))
	srv := serviceMock{}
	router := newRouter(&srv)
	response := responseWriterMock{}
	response.On("WriteHeader", http.StatusInternalServerError).Return(nil)
	response.On("Header").Return(http.Header{})
	response.On("Write", mocking.Anything).Return(0, nil)
	request := http.Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost",
			Path:   "/mock/mapping",
		},
		Method: http.MethodPost,
		Body:   body,
	}
	router.server.ServeHTTP(&response, &request)
	response.AssertCalled(t, "WriteHeader", http.StatusInternalServerError)
	response.AssertCalled(t, "Header")
	response.AssertCalled(t, "Write", mocking.Anything)
	response.AssertExpectations(t)
}

func TestMappingEntrypointMockNotFoundWhenCallService(t *testing.T) {
	body := io.NopCloser(strings.NewReader("{}"))
	srv := serviceMock{}
	router := newRouter(&srv)
	response := responseWriterMock{}
	response.On("WriteHeader", http.StatusNotFound).Return(nil)
	response.On("Header").Return(http.Header{})
	response.On("Write", mocking.Anything).Return(0, nil)
	srv.On("Add", mocking.AnythingOfType("mock.mockDTO")).Return(nil, mockNotFound(httpRequest{}))
	request := http.Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost",
			Path:   "/mock/mapping",
		},
		Method: http.MethodPost,
		Body:   body,
	}
	router.server.ServeHTTP(&response, &request)
	response.AssertCalled(t, "WriteHeader", http.StatusNotFound)
	response.AssertCalled(t, "Header")
	response.AssertCalled(t, "Write", mocking.Anything)
	response.AssertExpectations(t)
}

func TestMappingEntrypointInvalidMockWhenCallService(t *testing.T) {
	body := io.NopCloser(strings.NewReader("{}"))
	srv := serviceMock{}
	router := newRouter(&srv)
	response := responseWriterMock{}
	response.On("WriteHeader", http.StatusBadRequest).Return(nil)
	response.On("Header").Return(http.Header{})
	response.On("Write", mocking.Anything).Return(0, nil)
	srv.On("Add", mocking.AnythingOfType("mock.mockDTO")).Return(nil, invalidRequest("any cause"))
	request := http.Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost",
			Path:   "/mock/mapping",
		},
		Method: http.MethodPost,
		Body:   body,
	}
	router.server.ServeHTTP(&response, &request)
	response.AssertCalled(t, "WriteHeader", http.StatusBadRequest)
	response.AssertCalled(t, "Header")
	response.AssertCalled(t, "Write", mocking.Anything)
	response.AssertExpectations(t)
}

func TestMappingEntrypointOK(t *testing.T) {
	body := io.NopCloser(strings.NewReader("{}"))
	srv := serviceMock{}
	router := newRouter(&srv)
	response := responseWriterMock{}
	response.On("WriteHeader", http.StatusOK).Return(nil)
	response.On("Header").Return(http.Header{})
	response.On("Write", mocking.Anything).Return(0, nil)
	srv.On("Add", mocking.AnythingOfType("mock.mockDTO")).Return(&addMockResponse{}, nil)
	request := http.Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost",
			Path:   "/mock/mapping",
		},
		Method: http.MethodPost,
		Body:   body,
	}
	router.server.ServeHTTP(&response, &request)
	response.AssertCalled(t, "WriteHeader", http.StatusOK)
	response.AssertCalled(t, "Header")
	response.AssertCalled(t, "Write", mocking.Anything)
	response.AssertExpectations(t)
}

func TestMappingServeMockServiceError(t *testing.T) {
	srv := serviceMock{}
	router := newRouter(&srv)
	response := responseWriterMock{}
	response.On("WriteHeader", http.StatusInternalServerError).Return(nil)
	response.On("Header").Return(http.Header{})
	response.On("Write", mocking.Anything).Return(0, nil)
	srv.On("Match", mocking.AnythingOfType("mock.httpRequest")).Return(nil, Error{Code: "unknown"})
	request := http.Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost",
			Path:   "/test-url",
		},
		Method: http.MethodGet,
	}
	router.server.ServeHTTP(&response, &request)
	response.AssertCalled(t, "WriteHeader", http.StatusInternalServerError)
	response.AssertCalled(t, "Header")
	response.AssertCalled(t, "Write", mocking.Anything)
	response.AssertExpectations(t)
}

func TestMappingServeMockOK(t *testing.T) {
	srv := serviceMock{}
	router := newRouter(&srv)
	response := responseWriterMock{}
	response.On("WriteHeader", http.StatusAccepted).Return(nil)
	response.On("Header").Return(http.Header{})
	response.On("Write", mocking.Anything).Return(0, nil)
	srv.On("Match", mocking.AnythingOfType("mock.httpRequest")).Return(&httpResponse{Status: 202, Headers: map[string]string{"Content-Type": "application/json"}}, nil)
	request := http.Request{
		URL: &url.URL{
			Scheme:   "http",
			Host:     "localhost",
			Path:     "/test-url",
			RawQuery: "?app_name=topone",
		},
		Method: http.MethodGet,
		Header: map[string][]string{"Content-type": {"application/json"}},
	}
	router.server.ServeHTTP(&response, &request)
	response.AssertCalled(t, "WriteHeader", http.StatusAccepted)
	response.AssertCalled(t, "Header")
	response.AssertCalled(t, "Write", mocking.Anything)
	response.AssertExpectations(t)
}
