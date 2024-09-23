package mock

import (
	"net/http"

	mocking "github.com/stretchr/testify/mock"
)

type repositoryMock struct {
	mocking.Mock
}

func (r *repositoryMock) Save(info mock) error {
	args := r.Called(info)
	return args.Error(0)
}
func (r *repositoryMock) GetAll() []mock {
	args := r.Called()
	if args[0] == nil {
		return nil
	}
	return args[0].([]mock)
}

type serviceMock struct {
	mocking.Mock
}

func (r *serviceMock) Add(mock mockDTO) (*addMockResponse, error) {
	args := r.Called(mock)
	var r1 *addMockResponse
	if args.Get(0) != nil {
		r1 = args.Get(0).(*addMockResponse)
	}
	return r1, args.Error(1)
}

func (r *serviceMock) Match(request httpRequest) (*httpResponse, error) {
	args := r.Called(request)
	var r1 *httpResponse
	if args.Get(0) != nil {
		r1 = args.Get(0).(*httpResponse)
	}
	return r1, args.Error(1)
}

type responseWriterMock struct {
	mocking.Mock
}

func (resp *responseWriterMock) Header() http.Header {
	args := resp.Called()
	return args.Get(0).(http.Header)
}

func (resp *responseWriterMock) Write(data []byte) (int, error) {
	args := resp.Called(data)
	return args.Int(0), args.Error(1)
}

func (resp *responseWriterMock) WriteHeader(statusCode int) {
	resp.Called(statusCode)
}
