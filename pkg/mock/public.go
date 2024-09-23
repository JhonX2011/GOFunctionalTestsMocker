package mock

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const (
	operatorEqual    = "equal_to"
	operatorContains = "contains"
	operatorPattern  = "pattern"
)

type mockDTO struct {
	ID       string       `json:"id"`
	Request  *requestDTO  `json:"request"`
	Response *responseDTO `json:"response"`
}

type requestDTO struct {
	URL             map[string]string            `json:"url"`
	Method          *string                      `json:"method"`
	Headers         map[string]map[string]string `json:"headers"`
	QueryParameters map[string]map[string]string `json:"query_parameters"`
	Priority        int                          `json:"priority"`
	Body            map[string]string            `json:"body"`
}

type responseDTO struct {
	Status  int               `json:"status"`
	Body    json.RawMessage   `json:"body"`
	Headers map[string]string `json:"headers"`
}

type addMockResponse struct {
	ID string `json:"id"`
}

func (dto mockDTO) toAggregate() (*mock, error) {
	if dto.Request == nil {
		return nil, invalidRequest("the mock request could not be a null")
	}
	if dto.Response == nil {
		return nil, invalidRequest("the mock response could not be a null")
	}
	id := dto.ID
	if id == "" {
		uid, _ := uuid.NewUUID()
		id = uid.String()
	}
	request, err := toRequestMatch(dto)
	if err != nil {
		return nil, err
	}
	return &mock{
		ID:      id,
		Request: *request,
		Response: httpResponse{
			Status:  dto.Response.Status,
			Body:    dto.Response.Body,
			Headers: dto.Response.Headers,
		},
	}, nil
}

func toRequestMatch(dto mockDTO) (*requestMatch, error) {
	urlCondition, err := buildSimplexConditionFromMap(dto.Request.URL)
	if err != nil {
		return nil, err
	}
	headers, err := buildComplexCondition(dto.Request.Headers)
	if err != nil {
		return nil, err
	}
	queryParams, err := buildComplexCondition(dto.Request.QueryParameters)
	if err != nil {
		return nil, err
	}
	body, err := buildSimplexConditionFromMap(dto.Request.Body)
	if err != nil {
		return nil, err
	}

	return &requestMatch{
		URL:             urlCondition,
		Method:          dto.Request.Method,
		Headers:         headers,
		QueryParameters: queryParams,
		Priority:        dto.Request.Priority,
		Body:            body,
	}, nil
}

func buildComplexCondition(headers map[string]map[string]string) (complexConditions, error) {
	var conditionSlice complexConditions
	for field, rawCondition := range headers {
		rawOperator, value := getFirstMapEntry(rawCondition)
		if rawOperator == nil || value == nil {
			return nil, invalidRequest(fmt.Sprintf("the field %s has not any condition.", field))
		}
		op := fromString(*rawOperator)
		if op == undefined {
			return nil, invalidRequest(fmt.Sprintf("the operator %s is not supported.", *rawOperator))
		}
		condition := complexCondition{
			simplexCondition: simplexCondition{
				operator: op,
				value:    *value,
			},
			field: field,
		}
		conditionSlice = append(conditionSlice, condition)
	}
	return conditionSlice, nil
}

func buildSimplexConditionFromMap(data map[string]string) (*simplexCondition, error) {
	key, value := getFirstMapEntry(data)
	return buildSimpleCondition(key, value)
}

func buildSimpleCondition(key *string, value *string) (*simplexCondition, error) {
	if key == nil || value == nil {
		return nil, nil
	}
	op := fromString(*key)
	if op == undefined {
		return nil, invalidRequest(fmt.Sprintf("the operator %s is not supported.", *key))
	}
	return &simplexCondition{
		operator: op,
		value:    *value,
	}, nil
}

type requestBuilder struct {
	method          *string
	body            map[string]string
	url             map[string]string
	headers         map[string]map[string]string
	queryParameters map[string]map[string]string
	priority        int
}

type responseBuilder struct {
	status  int
	body    []byte
	headers map[string]string
}

func Request() RequestBuilder {
	return &requestBuilder{}
}

func Response() ResponseBuilder {
	return &responseBuilder{}
}

type RequestBuilder interface {
	URLEqualsTo(value string) RequestBuilder
	URLContains(value string) RequestBuilder
	URLPattern(value string) RequestBuilder
	Method(value string) RequestBuilder
	WithPriority(value int) RequestBuilder
	HeaderIsEqualTo(field string, value string) RequestBuilder
	HeaderContains(field string, value string) RequestBuilder
	HeaderPatternIs(field string, value string) RequestBuilder
	ParamIsEqualTo(field string, value string) RequestBuilder
	ParamContains(field string, value string) RequestBuilder
	ParamPatternIs(field string, value string) RequestBuilder
	BodyEqualsTo(body string) RequestBuilder
	BodyContains(part string) RequestBuilder
	BodyPatternIs(pattern string) RequestBuilder
	Build() *requestDTO
}

func (req *requestBuilder) addUrlEntry(key string, value string) RequestBuilder {
	if req.url == nil {
		req.url = map[string]string{}
	}
	req.url[key] = value
	return req
}
func (req *requestBuilder) URLEqualsTo(value string) RequestBuilder {
	return req.addUrlEntry(operatorEqual, value)
}
func (req *requestBuilder) URLContains(value string) RequestBuilder {
	return req.addUrlEntry(operatorContains, value)
}
func (req *requestBuilder) URLPattern(value string) RequestBuilder {
	return req.addUrlEntry(operatorPattern, value)
}
func (req *requestBuilder) Method(value string) RequestBuilder {
	req.method = &value
	return req
}
func (req *requestBuilder) WithPriority(value int) RequestBuilder {
	req.priority = value
	return req
}
func (req *requestBuilder) addHeader(field string, key string, value string) RequestBuilder {
	if req.headers == nil {
		req.headers = map[string]map[string]string{}
	}
	req.headers[field] = map[string]string{key: value}
	return req
}
func (req *requestBuilder) HeaderIsEqualTo(field string, value string) RequestBuilder {
	return req.addHeader(field, operatorEqual, value)
}
func (req *requestBuilder) HeaderContains(field string, value string) RequestBuilder {
	return req.addHeader(field, operatorContains, value)
}
func (req *requestBuilder) HeaderPatternIs(field string, value string) RequestBuilder {
	return req.addHeader(field, operatorPattern, value)
}

func (req *requestBuilder) addParam(field string, key string, value string) RequestBuilder {
	if req.queryParameters == nil {
		req.queryParameters = map[string]map[string]string{}
	}
	req.queryParameters[field] = map[string]string{key: value}
	return req
}

func (req *requestBuilder) ParamIsEqualTo(field string, value string) RequestBuilder {
	return req.addParam(field, operatorEqual, value)
}
func (req *requestBuilder) ParamContains(field string, value string) RequestBuilder {
	return req.addParam(field, operatorContains, value)
}
func (req *requestBuilder) ParamPatternIs(field string, value string) RequestBuilder {
	return req.addParam(field, operatorPattern, value)
}

func (req *requestBuilder) addBodyMatch(key string, value string) RequestBuilder {
	if req.body == nil {
		req.body = map[string]string{}
	}
	req.body[key] = value
	return req
}

func (req *requestBuilder) BodyEqualsTo(body string) RequestBuilder {
	return req.addBodyMatch(operatorEqual, body)
}
func (req *requestBuilder) BodyContains(part string) RequestBuilder {
	return req.addBodyMatch(operatorContains, part)
}
func (req *requestBuilder) BodyPatternIs(pattern string) RequestBuilder {
	return req.addBodyMatch(operatorPattern, pattern)
}

func (req *requestBuilder) Build() *requestDTO {
	return &requestDTO{
		URL:             req.url,
		Method:          req.method,
		Headers:         req.headers,
		QueryParameters: req.queryParameters,
		Priority:        req.priority,
		Body:            req.body,
	}
}

type ResponseBuilder interface {
	WithStatus(value int) ResponseBuilder
	WithBody(value []byte) ResponseBuilder
	WithBodyAsString(value string) ResponseBuilder
	WithHeader(name string, value string) ResponseBuilder
	WithHeaders(value map[string]string) ResponseBuilder
	Build() *responseDTO
}

func (res *responseBuilder) WithStatus(value int) ResponseBuilder {
	res.status = value
	return res
}
func (res *responseBuilder) WithBody(value []byte) ResponseBuilder {
	res.body = value
	return res
}
func (res *responseBuilder) WithBodyAsString(value string) ResponseBuilder {
	res.body = []byte(value)
	return res
}
func (res *responseBuilder) WithHeader(name string, value string) ResponseBuilder {
	if res.headers == nil {
		res.headers = map[string]string{}
	}
	res.headers[name] = value
	return res
}
func (res *responseBuilder) WithHeaders(value map[string]string) ResponseBuilder {
	res.headers = value
	return res
}
func (res *responseBuilder) Build() *responseDTO {
	return &responseDTO{
		Status:  res.status,
		Body:    res.body,
		Headers: res.headers,
	}
}

func getFirstMapEntry(m map[string]string) (*string, *string) {
	for key, value := range m {
		return &key, &value
	}
	return nil, nil
}

func fromString(name string) operator {
	switch name {
	case "equal_to":
		return equal
	case "contains":
		return contains
	case "pattern":
		return pattern
	default:
		return undefined
	}
}
