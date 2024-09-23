package mock

import "fmt"

type Error struct {
	Err         error  `json:"error"`
	Cause       string `json:"cause"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

const (
	errorTemplate      = "[Err: %v, Cause: %v, Code: %v, Description: %v]"
	invalidRequestCode = "invalid_request"
	mockNotFoundCode   = "mock_not_found"
)

func (err Error) Error() string {
	return fmt.Sprintf(errorTemplate, err.Err, err.Cause, err.Code, err.Description)
}

func invalidRequest(cause string) error {
	return Error{
		Code:        invalidRequestCode,
		Description: cause,
		Cause:       cause,
	}
}

func mockNotFound(request httpRequest) error {
	description := fmt.Sprintf("mapping not found for request %v.", request)
	return Error{
		Code:        mockNotFoundCode,
		Description: description,
		Cause:       description,
	}
}
