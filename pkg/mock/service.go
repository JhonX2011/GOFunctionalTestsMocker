package mock

import (
	"sort"
)

type Service interface {
	Add(mock mockDTO) (*addMockResponse, error)
	Match(request httpRequest) (*httpResponse, error)
}

type mockService struct {
	repository Repository
}

func newService(repository Repository) Service {
	return &mockService{
		repository: repository,
	}
}

func (instance *mockService) Add(mock mockDTO) (*addMockResponse, error) {
	err := validate(mock)
	if err != nil {
		LogInfo("error validating mock data")
		return nil, err
	}
	aggregate, err := mock.toAggregate()
	if err != nil {
		LogInfo("error when convert mock to aggregate")
		return nil, err
	}
	err = instance.repository.Save(*aggregate)
	if err != nil {
		LogInfo("error saving mock into repository")
		return nil, err
	}
	return &addMockResponse{
		ID: aggregate.ID,
	}, nil
}
func (instance *mockService) Match(request httpRequest) (*httpResponse, error) {
	aggregates := instance.repository.GetAll()
	if aggregates == nil || len(aggregates) < 1 {
		LogInfo("no aggregates found from repository")
		return nil, mockNotFound(request)
	}
	var filteredAggregates []mock
	for _, aggregate := range aggregates {
		if aggregate.Request.IsExpected(request) {
			filteredAggregates = append(filteredAggregates, aggregate)
		}
	}
	if len(filteredAggregates) < 1 {
		LogInfo("no aggregates found for request %v", request)
		return nil, mockNotFound(request)
	}
	sort.SliceStable(filteredAggregates, func(i, j int) bool {
		return filteredAggregates[i].Request.Priority > filteredAggregates[j].Request.Priority
	})
	LogInfo("filter aggregates are %v", filteredAggregates)
	return &filteredAggregates[0].Response, nil
}

func validate(m mockDTO) error {
	if m.Request == nil {
		return invalidRequest("the mock request could not be a null")
	}
	if m.Response == nil {
		return invalidRequest("the mock response could not be a null")
	}
	if m.Request.URL == nil && m.Request.Method == nil && m.Request.Headers == nil && m.Request.QueryParameters == nil {
		return invalidRequest("the request has no conditions")
	}
	if m.Response.Status == 0 {
		return invalidRequest("the response status is required")
	}
	return nil
}
