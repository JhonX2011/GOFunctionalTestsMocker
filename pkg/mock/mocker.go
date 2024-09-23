package mock

import "fmt"

type Mocker interface {
	When(req *requestDTO) Expect
}

type Expect interface {
	ThenReturn(resp *responseDTO) error
}

type mocker struct {
	service Service
}
type expect struct {
	req     *requestDTO
	service Service
}

func (m *mocker) When(req *requestDTO) Expect {
	return &expect{
		service: m.service,
		req:     req,
	}
}
func (exp *expect) ThenReturn(resp *responseDTO) error {
	if exp.req == nil {
		return fmt.Errorf("the request builder expected could not be nil")
	}
	if resp == nil {
		return fmt.Errorf("the response builder could not be nil")
	}
	mock := mockDTO{
		Request:  exp.req,
		Response: resp,
	}
	_, err := exp.service.Add(mock)
	exp.req = nil
	exp.service = nil
	return err
}
