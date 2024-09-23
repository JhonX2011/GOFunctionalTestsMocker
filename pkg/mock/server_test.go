package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mocking "github.com/stretchr/testify/mock"
)

func TestNewServer(t *testing.T) {
	router, mocker := New()
	assert.NotNil(t, router)
	assert.NotNil(t, mocker)
}

func TestNewServerPanic(t *testing.T) {
	assert.NotPanics(t, func() { New() })
}

func TestMockRequestSuccess(t *testing.T) {
	srvMock := serviceMock{}
	mocker := internalNew(&srvMock)
	resp := &addMockResponse{}
	srvMock.On("Add", mocking.AnythingOfType("mock.mockDTO")).Return(resp, nil)
	err := mocker.When(
		Request().
			WithPriority(1).
			URLEqualsTo("/inventories/123123").
			Build(),
	).ThenReturn(
		Response().
			WithStatus(200).
			WithBodyAsString(`{"name":"pedro"}`).
			Build(),
	)
	assert.Nil(t, err)
	srvMock.AssertExpectations(t)
}

func TestMockRequestError(t *testing.T) {
	srvMock := serviceMock{}
	mocker := internalNew(&srvMock)
	srvMock.On("Add", mocking.AnythingOfType("mock.mockDTO")).Return(nil, invalidRequest("invalid"))
	err := mocker.When(
		Request().
			WithPriority(1).
			URLEqualsTo("/inventories/123123").
			Build(),
	).ThenReturn(
		Response().
			WithStatus(200).
			WithBodyAsString(`{"name":"pedro"}`).
			Build(),
	)
	assert.Error(t, err)
	srvMock.AssertExpectations(t)
}

func TestMockRequestWhenResponseBuilderIsNil(t *testing.T) {
	srvMock := serviceMock{}
	mocker := internalNew(&srvMock)
	err := mocker.When(
		Request().
			WithPriority(1).
			URLEqualsTo("/inventories/123123").
			Build(),
	).ThenReturn(nil)
	assert.Error(t, err)
	assert.Equal(t, "the response builder could not be nil", err.Error())
	srvMock.AssertExpectations(t)
	srvMock.AssertNotCalled(t, "Add")
}

func TestMockRequestWhenRequestBuilderIsNil(t *testing.T) {
	srvMock := serviceMock{}
	mocker := internalNew(&srvMock)
	err := mocker.When(nil).
		ThenReturn(
			Response().
				WithStatus(200).
				WithBodyAsString(`{"name":"pedro"}`).
				Build(),
		)
	assert.Error(t, err)
	assert.Equal(t, "the request builder expected could not be nil", err.Error())
	srvMock.AssertExpectations(t)
	srvMock.AssertNotCalled(t, "Add")
}
