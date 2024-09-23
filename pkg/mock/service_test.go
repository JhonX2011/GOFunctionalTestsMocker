package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mocking "github.com/stretchr/testify/mock"
)

func TestAddSuccess(t *testing.T) {
	method := "PUT"
	id := "21324123"
	url := "/any-url"
	responseBody := []byte(`{"results": 12312}`)
	responseStatus := 200
	responseHeaders := map[string]string{"Content-Type": "application/json"}
	m := mockDTO{
		ID: id,
		Request: &requestDTO{
			URL: map[string]string{
				"equal_to": url,
			},
			Method: &method,
			Headers: map[string]map[string]string{
				"Content-Type": {"contains": "application/json"},
			},
			QueryParameters: map[string]map[string]string{
				"query": {"pattern": "sql"},
			},
			Priority: 1,
		},
		Response: &responseDTO{
			Status:  responseStatus,
			Body:    responseBody,
			Headers: responseHeaders,
		},
	}
	repo := repositoryMock{}
	repo.On("Save", mocking.AnythingOfType("mock.mock")).Return(nil)
	service := newService(&repo)
	res, err := service.Add(m)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.ID)
	repo.AssertExpectations(t)
}

func TestAddRepositoryError(t *testing.T) {
	method := "PUT"
	id := "21324123"
	url := "/any-url"
	responseBody := []byte(`{"results": 12312}`)
	responseStatus := 200
	responseHeaders := map[string]string{"Content-Type": "application/json"}
	m := mockDTO{
		ID: id,
		Request: &requestDTO{
			URL: map[string]string{
				"equal_to": url,
			},
			Method: &method,
			Headers: map[string]map[string]string{
				"Content-Type": {"contains": "application/json"},
			},
			QueryParameters: map[string]map[string]string{
				"query": {"pattern": "sql"},
			},
			Priority: 1,
		},
		Response: &responseDTO{
			Status:  responseStatus,
			Body:    responseBody,
			Headers: responseHeaders,
		},
	}
	repo := repositoryMock{}
	repo.On("Save", mocking.AnythingOfType("mock.mock")).Return(invalidRequest("any cause"))
	service := newService(&repo)
	_, err := service.Add(m)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestAddInvalidRequestInfo(t *testing.T) {
	id := "21324123"
	url := "/any-url"
	responseBody := []byte(`{"results": 12312}`)
	responseStatus := 200
	responseHeaders := map[string]string{"Content-Type": "application/json"}
	m := mockDTO{
		ID: id,
		Request: &requestDTO{
			URL: map[string]string{
				"equal": url,
			},
			Priority: 1,
		},
		Response: &responseDTO{
			Status:  responseStatus,
			Body:    responseBody,
			Headers: responseHeaders,
		},
	}
	repo := repositoryMock{}
	service := newService(&repo)
	_, err := service.Add(m)
	assert.Error(t, err)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Add")
}

func TestAddInvalidMockRequest(t *testing.T) {
	responseBody := []byte(`{"results": 12312}`)
	responseStatus := 200
	responseHeaders := map[string]string{"Content-Type": "application/json"}
	m := mockDTO{
		Request: nil,
		Response: &responseDTO{
			Status:  responseStatus,
			Body:    responseBody,
			Headers: responseHeaders,
		},
	}
	repo := repositoryMock{}
	service := newService(&repo)
	_, err := service.Add(m)
	assert.Error(t, err)
	assert.Equal(t, "the mock request could not be a null", err.(Error).Cause)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Add")
}

func TestAddInvalidMockResponse(t *testing.T) {
	m := mockDTO{
		Request: &requestDTO{
			URL: map[string]string{
				"equal_to": "/test",
			},
		},
		Response: nil,
	}
	repo := repositoryMock{}
	service := newService(&repo)
	_, err := service.Add(m)
	assert.Error(t, err)
	assert.Equal(t, "the mock response could not be a null", err.(Error).Cause)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Add")
}

func TestAddInvalidMockResponseCode(t *testing.T) {
	m := mockDTO{
		Request: &requestDTO{
			URL: map[string]string{
				"equal_to": "/test",
			},
		},
		Response: &responseDTO{},
	}
	repo := repositoryMock{}
	service := newService(&repo)
	_, err := service.Add(m)
	assert.Error(t, err)
	assert.Equal(t, "the response status is required", err.(Error).Cause)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Add")
}

func TestAddInvalidMockRequestData(t *testing.T) {
	m := mockDTO{
		Request:  &requestDTO{},
		Response: &responseDTO{},
	}
	repo := repositoryMock{}
	service := newService(&repo)
	_, err := service.Add(m)
	assert.Error(t, err)
	assert.Equal(t, "the request has no conditions", err.(Error).Cause)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Add")
}

func TestMatchSuccess(t *testing.T) {
	aggregates := []mock{
		{
			ID: "1",
			Request: requestMatch{
				URL: &simplexCondition{
					operator: equal,
					value:    "/test",
				},
			},
			Response: httpResponse{
				Status: 404,
				Body:   []byte(`{"name": "any-name"}`),
			},
		},
		{
			ID: "2",
			Request: requestMatch{
				URL: &simplexCondition{
					operator: equal,
					value:    "/not-test",
				},
			},
			Response: httpResponse{
				Status: 200,
				Body:   []byte(`{"name": "other-name"}`),
			},
		},
	}
	req := httpRequest{
		URL: "/test",
	}
	repo := repositoryMock{}
	service := newService(&repo)
	repo.On("GetAll").Return(aggregates)
	resp, err := service.Match(req)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"name": "any-name"}`), resp.Body)
	assert.Equal(t, 404, resp.Status)
	repo.AssertExpectations(t)
}

func TestMatchSuccessOrderPriority(t *testing.T) {
	aggregates := []mock{
		{
			ID: "1",
			Request: requestMatch{
				URL: &simplexCondition{
					operator: equal,
					value:    "/test",
				},
				Priority: 100,
			},
			Response: httpResponse{
				Status: 404,
				Body:   []byte(`{"name": "any-name"}`),
			},
		},
		{
			ID: "2",
			Request: requestMatch{
				URL: &simplexCondition{
					operator: contains,
					value:    "test",
				},
				Priority: 150,
			},
			Response: httpResponse{
				Status: 200,
				Body:   []byte(`{"name": "other-name"}`),
			},
		},
	}
	req := httpRequest{
		URL: "/test",
	}
	repo := repositoryMock{}
	service := newService(&repo)
	repo.On("GetAll").Return(aggregates)
	resp, err := service.Match(req)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"name": "other-name"}`), resp.Body)
	assert.Equal(t, 200, resp.Status)
	repo.AssertExpectations(t)
}

func TestMatchWhenAllAggregatesAllFiltered(t *testing.T) {
	aggregates := []mock{
		{
			ID: "1",
			Request: requestMatch{
				URL: &simplexCondition{
					operator: equal,
					value:    "/test",
				},
			},
			Response: httpResponse{
				Status: 404,
				Body:   []byte(`{"name": "any-name"}`),
			},
		},
		{
			ID: "2",
			Request: requestMatch{
				URL: &simplexCondition{
					operator: equal,
					value:    "/not-test",
				},
			},
			Response: httpResponse{
				Status: 200,
				Body:   []byte(`{"name": "other-name"}`),
			},
		},
	}
	req := httpRequest{
		URL: "/other",
	}
	repo := repositoryMock{}
	service := newService(&repo)
	repo.On("GetAll").Return(aggregates)
	_, err := service.Match(req)
	assert.Error(t, err)
	assert.Equal(t, "mock_not_found", err.(Error).Code)
	repo.AssertExpectations(t)
}

func TestMatchNullAggregates(t *testing.T) {
	var aggregates []mock
	req := httpRequest{
		URL: "/test",
	}
	repo := repositoryMock{}
	service := newService(&repo)
	repo.On("GetAll").Return(aggregates)
	_, err := service.Match(req)
	assert.Error(t, err)
	assert.Equal(t, "mock_not_found", err.(Error).Code)
	repo.AssertExpectations(t)
}
