package mock

import (
	"regexp"
	"strings"
)

const (
	undefined operator = iota
	equal
	contains
	pattern
)

type biPredicate[T comparable, Z comparable] func(T, Z) bool

type complexConditions []complexCondition

type operator uint8

type simplexCondition struct {
	operator operator
	value    string
}

type complexCondition struct {
	simplexCondition
	field string
}

type mock struct {
	ID       string       `json:"id"`
	Request  requestMatch `json:"request"`
	Response httpResponse `json:"response"`
}

type requestMatch struct {
	URL             *simplexCondition `json:"url"`
	Method          *string           `json:"method"`
	Headers         complexConditions `json:"headers"`
	QueryParameters complexConditions `json:"query_parameters"`
	Body            *simplexCondition `json:"body"`
	Priority        int               `json:"priority"`
}

type httpRequest struct {
	URL             string
	Method          string
	Headers         map[string]string
	QueryParameters map[string]string
	Body            []byte
}

type httpResponse struct {
	Status  int               `json:"status"`
	Body    []byte            `json:"body"`
	Headers map[string]string `json:"headers"`
}

func equalsPredicate(value string, toCompare string) bool {
	return value == toCompare
}

func regexPredicate(value string, toCompare string) bool {
	reg, _ := regexp.Compile(value)
	return reg.MatchString(toCompare)
}

func containsPredicate(value string, toCompare string) bool {
	return strings.Contains(toCompare, value)
}

func (match *requestMatch) IsExpected(request httpRequest) bool {
	urlMatch := true
	if match.URL != nil {
		urlMatch = match.URL.test(request.URL)
	}
	methodMatch := true
	if match.Method != nil {
		methodMatch = *match.Method == request.Method
	}
	headerMatch := true
	if match.Headers != nil {
		headerMatch = match.Headers.match(request.Headers)
	}
	queryMatch := true
	if match.QueryParameters != nil {
		queryMatch = match.QueryParameters.match(request.QueryParameters)
	}
	bodyMatch := true
	if match.Body != nil {
		bodyMatch = match.Body.test(string(request.Body))
	}
	return urlMatch && methodMatch && headerMatch && queryMatch && bodyMatch
}

func (conditions complexConditions) match(params map[string]string) bool {
	for _, condition := range conditions {
		if !condition.test(params) {
			return false
		}
	}
	return true
}

func (c simplexCondition) test(value string) bool {
	predicate := c.operator.getPredicate()
	if predicate == nil {
		return false
	}
	return predicate(c.value, value)
}

func (c complexCondition) test(value map[string]string) bool {
	if v, exists := value[c.field]; exists {
		predicate := c.operator.getPredicate()
		if predicate == nil {
			return false
		}
		return predicate(c.value, v)
	}
	return false
}

func (o operator) getPredicate() biPredicate[string, string] {
	switch o {
	case contains:
		return containsPredicate
	case pattern:
		return regexPredicate
	case equal:
		return equalsPredicate
	default:
		return nil

	}
}
