package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidRequest(t *testing.T) {
	err := invalidRequest("any cause")
	assert.Equal(t, "[Err: <nil>, Cause: any cause, Code: invalid_request, Description: any cause]", err.Error())
}

func TestMockNotFound(t *testing.T) {
	err := mockNotFound(httpRequest{})
	assert.Equal(t, "[Err: <nil>, Cause: mapping not found for request {  map[] map[] []}., Code: mock_not_found, Description: mapping not found for request {  map[] map[] []}.]", err.Error())
}
