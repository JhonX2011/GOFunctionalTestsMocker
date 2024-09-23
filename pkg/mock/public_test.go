package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToAggregate(t *testing.T) {
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
			Body:     map[string]string{"contains": "any-body"},
			Priority: 1,
		},
		Response: &responseDTO{
			Status:  responseStatus,
			Body:    responseBody,
			Headers: responseHeaders,
		},
	}
	aggregate, err := m.toAggregate()
	assert.Nil(t, err)
	assert.NotNil(t, aggregate)
	assert.Equal(t, "21324123", aggregate.ID)
	req := aggregate.Request
	assert.Equal(t, 1, req.Priority)
	// url
	assert.NotNil(t, req)
	assert.Equal(t, method, *req.Method)
	assert.Equal(t, url, req.URL.value)
	assert.Equal(t, equal, req.URL.operator)
	// Headers
	assert.NotNil(t, req.Headers)
	assert.NotEmpty(t, req.Headers)
	assert.Equal(t, "Content-Type", req.Headers[0].field)
	assert.Equal(t, "application/json", req.Headers[0].value)
	assert.Equal(t, contains, req.Headers[0].operator)
	// Params
	assert.NotNil(t, req.QueryParameters)
	assert.NotEmpty(t, req.QueryParameters)
	assert.Equal(t, "query", req.QueryParameters[0].field)
	assert.Equal(t, "sql", req.QueryParameters[0].value)
	assert.Equal(t, pattern, req.QueryParameters[0].operator)
	// response
	assert.Equal(t, responseHeaders, aggregate.Response.Headers)
	assert.Equal(t, responseBody, aggregate.Response.Body)
	assert.Equal(t, responseStatus, aggregate.Response.Status)
}

func TestToAggregateInvalidUrlOperator(t *testing.T) {
	url := "/any-url"
	responseBody := []byte(`{"results": 12312}`)
	responseStatus := 200
	responseHeaders := map[string]string{"Content-Type": "application/json"}
	m := mockDTO{
		Request: &requestDTO{
			URL: map[string]string{
				"equals": url,
			},
		},
		Response: &responseDTO{
			Status:  responseStatus,
			Body:    responseBody,
			Headers: responseHeaders,
		},
	}
	_, err := m.toAggregate()
	assert.Error(t, err)
	assert.Equal(t, "invalid_request", err.(Error).Code)
	assert.Equal(t, "the operator equals is not supported.", err.(Error).Description)
}

func TestToAggregateInvalidHeaderOperator(t *testing.T) {
	responseBody := []byte(`{"results": 12312}`)
	responseStatus := 200
	responseHeaders := map[string]string{"Content-Type": "application/json"}
	m := mockDTO{
		Request: &requestDTO{
			Headers: map[string]map[string]string{
				"Accept-Encoding": {"match": "gzip"},
			},
		},
		Response: &responseDTO{
			Status:  responseStatus,
			Body:    responseBody,
			Headers: responseHeaders,
		},
	}
	_, err := m.toAggregate()
	assert.Error(t, err)
	assert.Equal(t, "invalid_request", err.(Error).Code)
	assert.Equal(t, "the operator match is not supported.", err.(Error).Description)
}

func TestToAggregateInvalidQueryParamInvalidOperator(t *testing.T) {
	responseBody := []byte(`{"results": 12312}`)
	responseStatus := 200
	responseHeaders := map[string]string{"Content-Type": "application/json"}
	m := mockDTO{
		Request: &requestDTO{
			QueryParameters: map[string]map[string]string{
				"version": {"match": "1.0.0"},
			},
		},
		Response: &responseDTO{
			Status:  responseStatus,
			Body:    responseBody,
			Headers: responseHeaders,
		},
	}
	_, err := m.toAggregate()
	assert.Error(t, err)
	assert.Equal(t, "invalid_request", err.(Error).Code)
	assert.Equal(t, "the operator match is not supported.", err.(Error).Description)
}

func TestToAggregateInvalidBodyOperator(t *testing.T) {
	responseBody := []byte(`{"results": 12312}`)
	responseStatus := 200
	responseHeaders := map[string]string{"Content-Type": "application/json"}
	m := mockDTO{
		Request: &requestDTO{
			Body: map[string]string{"invalid-condition": "any-value"},
		},
		Response: &responseDTO{
			Status:  responseStatus,
			Body:    responseBody,
			Headers: responseHeaders,
		},
	}
	_, err := m.toAggregate()
	assert.Error(t, err)
	assert.Equal(t, "invalid_request", err.(Error).Code)
	assert.Equal(t, "the operator invalid-condition is not supported.", err.(Error).Description)
}

func TestToAggregateQueryParamWithEmptyCondition(t *testing.T) {
	responseBody := []byte(`{"results": 12312}`)
	responseStatus := 200
	responseHeaders := map[string]string{"Content-Type": "application/json"}
	m := mockDTO{
		Request: &requestDTO{
			QueryParameters: map[string]map[string]string{
				"version": nil,
			},
		},
		Response: &responseDTO{
			Status:  responseStatus,
			Body:    responseBody,
			Headers: responseHeaders,
		},
	}
	_, err := m.toAggregate()
	assert.Error(t, err)
	assert.Equal(t, "invalid_request", err.(Error).Code)
	assert.Equal(t, "the field version has not any condition.", err.(Error).Description)
}

func TestToAggregateHeadersWithEmptyCondition(t *testing.T) {
	responseBody := []byte(`{"results": 12312}`)
	responseStatus := 200
	responseHeaders := map[string]string{"Content-Type": "application/json"}
	m := mockDTO{
		Request: &requestDTO{
			Headers: map[string]map[string]string{
				"Accept-Version": nil,
			},
		},
		Response: &responseDTO{
			Status:  responseStatus,
			Body:    responseBody,
			Headers: responseHeaders,
		},
	}
	_, err := m.toAggregate()
	assert.Error(t, err)
	assert.Equal(t, "invalid_request", err.(Error).Code)
	assert.Equal(t, "the field Accept-Version has not any condition.", err.(Error).Description)
}

func TestToAggregateRequestNil(t *testing.T) {
	responseBody := []byte(`{"results": 12312}`)
	responseStatus := 200
	responseHeaders := map[string]string{"Content-Type": "application/json"}
	m := mockDTO{
		Response: &responseDTO{
			Status:  responseStatus,
			Body:    responseBody,
			Headers: responseHeaders,
		},
	}
	_, err := m.toAggregate()
	assert.Error(t, err)
	assert.Equal(t, "invalid_request", err.(Error).Code)
	assert.Equal(t, "the mock request could not be a null", err.(Error).Description)

}

func TestToAggregateResponseNil(t *testing.T) {
	method := "PUT"
	m := mockDTO{
		Request: &requestDTO{
			Method: &method,
		},
	}
	_, err := m.toAggregate()
	assert.Error(t, err)
	assert.Equal(t, "invalid_request", err.(Error).Code)
	assert.Equal(t, "the mock response could not be a null", err.(Error).Description)
}

func TestRequestBuilderWithContainsCondition(t *testing.T) {
	url := "/test"
	priority := 1
	method := "GET"
	req := Request().
		Method(method).
		HeaderContains("Accept-Encoding", "gzip").
		HeaderContains("Content-Type", "application/json").
		ParamContains("version", "0.1").
		ParamContains("type", "any-type").
		BodyContains("body-part").
		URLContains(url).
		WithPriority(priority).
		Build()
	assert.Equal(t, method, *req.Method)
	assert.Equal(t, priority, req.Priority)
	assert.Equal(t, map[string]string{"contains": url}, req.URL)
	assert.Equal(t, map[string]map[string]string{"Accept-Encoding": {"contains": "gzip"}, "Content-Type": {"contains": "application/json"}}, req.Headers)
	assert.Equal(t, map[string]map[string]string{"version": {"contains": "0.1"}, "type": {"contains": "any-type"}}, req.QueryParameters)
	assert.Equal(t, map[string]string{"contains": "body-part"}, req.Body)
}

func TestRequestBuilderWithEqualsCondition(t *testing.T) {
	url := "/test"
	priority := 1
	method := "GET"
	req := Request().
		Method(method).
		HeaderIsEqualTo("Accept-Encoding", "gzip").
		HeaderIsEqualTo("Content-Type", "application/json").
		ParamIsEqualTo("version", "0.1").
		ParamIsEqualTo("type", "any-type").
		BodyEqualsTo("any-condition").
		URLEqualsTo(url).
		WithPriority(priority).
		Build()
	assert.Equal(t, method, *req.Method)
	assert.Equal(t, priority, req.Priority)
	assert.Equal(t, map[string]string{"equal_to": url}, req.URL)
	assert.Equal(t, map[string]map[string]string{"Accept-Encoding": {"equal_to": "gzip"}, "Content-Type": {"equal_to": "application/json"}}, req.Headers)
	assert.Equal(t, map[string]map[string]string{"version": {"equal_to": "0.1"}, "type": {"equal_to": "any-type"}}, req.QueryParameters)
	assert.Equal(t, 1, len(req.Body))
	assert.Equal(t, map[string]string{"equal_to": "any-condition"}, req.Body)
}

func TestRequestBuilderWithPatternCondition(t *testing.T) {
	url := "/test"
	priority := 1
	method := "GET"
	req := Request().
		Method(method).
		HeaderPatternIs("Accept-Encoding", "gzip").
		HeaderPatternIs("Content-Type", "application/json").
		ParamPatternIs("version", "0.1").
		ParamPatternIs("type", "any-type").
		BodyPatternIs("any-pattern").
		URLPattern(url).
		WithPriority(priority).
		Build()
	assert.Equal(t, method, *req.Method)
	assert.Equal(t, priority, req.Priority)
	assert.Equal(t, map[string]string{"pattern": url}, req.URL)
	assert.Equal(t, map[string]map[string]string{"Accept-Encoding": {"pattern": "gzip"}, "Content-Type": {"pattern": "application/json"}}, req.Headers)
	assert.Equal(t, map[string]map[string]string{"version": {"pattern": "0.1"}, "type": {"pattern": "any-type"}}, req.QueryParameters)
	assert.Equal(t, map[string]string{"pattern": "any-pattern"}, req.Body)
}

func TestResponseBuilderWithBodyAsString(t *testing.T) {
	status := 200
	body := `{"name": "any-name"}`
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	resp := Response().
		WithStatus(status).
		WithBodyAsString(body).
		WithHeaders(headers).
		WithHeader("Accept", "application/json").
		Build()
	assert.Equal(t, body, string(resp.Body))
	assert.Equal(t, status, resp.Status)
	assert.Equal(t, map[string]string{"Content-Type": "application/json", "Accept": "application/json"}, resp.Headers)
}

func TestResponseBuilderWithBodyAsByteArray(t *testing.T) {
	status := 200
	body := []byte(`{"name": "any-name"}`)
	resp := Response().
		WithStatus(status).
		WithBody(body).
		WithHeader("Accept", "application/json").
		WithHeader("Content-Type", "application/json").
		Build()
	assert.Equal(t, string(body), string(resp.Body))
	assert.Equal(t, status, resp.Status)
	assert.Equal(t, map[string]string{"Content-Type": "application/json", "Accept": "application/json"}, resp.Headers)
}
