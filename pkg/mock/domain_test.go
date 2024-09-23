package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	postMethod = "POST"
	getMethod  = "GET"
)

func TestUrlEquals(t *testing.T) {
	reqMatch := requestMatch{
		URL: &simplexCondition{
			operator: equal,
			value:    "/test",
		},
	}
	req := httpRequest{
		URL: "/test",
	}
	expected := reqMatch.IsExpected(req)
	assert.True(t, expected)
}

func TestUrlNotEquals(t *testing.T) {
	reqMatch := requestMatch{
		URL: &simplexCondition{
			operator: equal,
			value:    "/test",
		},
	}
	req := httpRequest{
		URL: "/testing",
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestUrlContains(t *testing.T) {
	reqMatch := requestMatch{
		URL: &simplexCondition{
			operator: contains,
			value:    "/test",
		},
	}
	req := httpRequest{
		URL: "/test/123143",
	}
	expected := reqMatch.IsExpected(req)
	assert.True(t, expected)
}
func TestUrlNotContains(t *testing.T) {
	reqMatch := requestMatch{
		URL: &simplexCondition{
			operator: contains,
			value:    "/testing",
		},
	}
	req := httpRequest{
		URL: "/test/123143",
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestUrlPatternMatch(t *testing.T) {
	reqMatch := requestMatch{
		URL: &simplexCondition{
			operator: pattern,
			value:    "^/test.*",
		},
	}
	req := httpRequest{
		URL: "/test?name=any",
	}
	expected := reqMatch.IsExpected(req)
	assert.True(t, expected)
}

func TestUrlPatternNotMatch(t *testing.T) {
	reqMatch := requestMatch{
		URL: &simplexCondition{
			operator: pattern,
			value:    "^/test.*",
		},
	}
	req := httpRequest{
		URL: "/fbm/test?name=any",
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestOperatorUndefined(t *testing.T) {
	reqMatch := requestMatch{
		URL: &simplexCondition{
			operator: undefined,
			value:    "/test",
		},
	}
	req := httpRequest{
		URL: "/test",
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestMethodEquals(t *testing.T) {
	matchMethod := getMethod
	reqMatch := requestMatch{
		Method: &matchMethod,
	}
	req := httpRequest{
		Method: getMethod,
	}
	expected := reqMatch.IsExpected(req)
	assert.True(t, expected)
}

func TestMethodNotEquals(t *testing.T) {
	matchMethod := getMethod
	reqMatch := requestMatch{
		Method: &matchMethod,
	}
	req := httpRequest{
		Method: postMethod,
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestUrlAndMethodEquals(t *testing.T) {
	matchMethod := getMethod
	url := "/test"
	reqMatch := requestMatch{
		URL: &simplexCondition{
			operator: equal,
			value:    url,
		},
		Method: &matchMethod,
	}
	req := httpRequest{
		Method: getMethod,
		URL:    url,
	}
	expected := reqMatch.IsExpected(req)
	assert.True(t, expected)
}

func TestUrlEqualsAndMethodNot(t *testing.T) {
	matchMethod := getMethod
	url := "/test"
	reqMatch := requestMatch{
		URL: &simplexCondition{
			operator: equal,
			value:    url,
		},
		Method: &matchMethod,
	}
	req := httpRequest{
		Method: postMethod,
		URL:    url,
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestHeaderEquals(t *testing.T) {
	reqMatch := requestMatch{
		Headers: complexConditions{
			complexCondition{
				simplexCondition: simplexCondition{
					operator: equal,
					value:    "application/json",
				},
				field: "Content-Type",
			},
			complexCondition{
				simplexCondition: simplexCondition{
					operator: contains,
					value:    "gzip",
				},
				field: "Accept-Encoding",
			},
			complexCondition{
				simplexCondition: simplexCondition{
					operator: pattern,
					value:    "^Chrome.*",
				},
				field: "User-Agent",
			},
		},
	}
	req := httpRequest{
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"Accept-Encoding": "gzip, deflate, br",
			"User-Agent":      "Chrome/Chromuim v1.2.3",
		},
	}
	expected := reqMatch.IsExpected(req)
	assert.True(t, expected)
}

func TestHeaderNotEquals(t *testing.T) {
	reqMatch := requestMatch{
		Headers: complexConditions{
			complexCondition{
				simplexCondition: simplexCondition{
					operator: equal,
					value:    "application/json",
				},
				field: "Content-Type",
			},
		},
	}
	req := httpRequest{
		Headers: map[string]string{
			"Content-Type":    "application/xml",
			"Accept-Encoding": "gzip, deflate, br",
			"User-Agent":      "Chrome/Chromuim v1.2.3",
		},
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestHeaderNotContains(t *testing.T) {
	reqMatch := requestMatch{
		Headers: complexConditions{
			complexCondition{
				simplexCondition: simplexCondition{
					operator: contains,
					value:    "gzip",
				},
				field: "Accept-Encoding",
			},
		},
	}
	req := httpRequest{
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"Accept-Encoding": "zip, deflate, br",
			"User-Agent":      "Chrome/Chromuim v1.2.3",
		},
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestHeaderPatternNotMatch(t *testing.T) {
	reqMatch := requestMatch{
		Headers: complexConditions{
			complexCondition{
				simplexCondition: simplexCondition{
					operator: pattern,
					value:    "^Chrome.*",
				},
				field: "User-Agent",
			},
		},
	}
	req := httpRequest{
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"Accept-Encoding": "zip, deflate, br",
			"User-Agent":      "Firefox v1.2.3",
		},
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestHeaderUndefinedOperator(t *testing.T) {
	reqMatch := requestMatch{
		Headers: complexConditions{
			complexCondition{
				simplexCondition: simplexCondition{
					operator: undefined,
					value:    "^Firefox.*",
				},
				field: "User-Agent",
			},
		},
	}
	req := httpRequest{
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"Accept-Encoding": "zip, deflate, br",
			"User-Agent":      "Firefox v1.2.3",
		},
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestWhenAHeaderNotExists(t *testing.T) {
	reqMatch := requestMatch{
		Headers: complexConditions{
			complexCondition{
				simplexCondition: simplexCondition{
					operator: undefined,
					value:    "^Firefox.*",
				},
				field: "User-Agent",
			},
		},
	}
	req := httpRequest{
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"Accept-Encoding": "zip, deflate, br",
		},
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestQueryParamsEquals(t *testing.T) {
	reqMatch := requestMatch{
		QueryParameters: complexConditions{
			complexCondition{
				simplexCondition: simplexCondition{
					operator: equal,
					value:    "MLB",
				},
				field: "site",
			},
			complexCondition{
				simplexCondition: simplexCondition{
					operator: contains,
					value:    "heavy",
				},
				field: "categories",
			},
			complexCondition{
				simplexCondition: simplexCondition{
					operator: pattern,
					value:    "^solid*",
				},
				field: "query",
			},
		},
	}
	req := httpRequest{
		QueryParameters: map[string]string{
			"site":       "MLB",
			"categories": "dry,heavy,child",
			"query":      "solid description for that",
		},
	}
	expected := reqMatch.IsExpected(req)
	assert.True(t, expected)
}

func TestQueryParamsNotEquals(t *testing.T) {
	reqMatch := requestMatch{
		QueryParameters: complexConditions{
			complexCondition{
				simplexCondition: simplexCondition{
					operator: equal,
					value:    "MLB",
				},
				field: "site",
			},
		},
	}
	req := httpRequest{
		QueryParameters: map[string]string{
			"site":       "MLM",
			"categories": "dry,heavy,child",
			"query":      "solid description for that",
		},
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestQueryParamsNotContains(t *testing.T) {
	reqMatch := requestMatch{
		QueryParameters: complexConditions{
			complexCondition{
				simplexCondition: simplexCondition{
					operator: contains,
					value:    "heavy",
				},
				field: "categories",
			},
		},
	}
	req := httpRequest{
		QueryParameters: map[string]string{
			"site":       "MLB",
			"categories": "dry,child",
			"query":      "solid description for that",
		},
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestQueryParamsPatternNotMatch(t *testing.T) {
	reqMatch := requestMatch{
		QueryParameters: complexConditions{
			complexCondition{
				simplexCondition: simplexCondition{
					operator: pattern,
					value:    "^solid*",
				},
				field: "query",
			},
		},
	}
	req := httpRequest{
		QueryParameters: map[string]string{
			"site":       "MLB",
			"categories": "dry,heavy,child",
			"query":      "not solid description for that",
		},
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestHeaderAndQueryEqualsAndUrlNot(t *testing.T) {
	reqMatch := requestMatch{
		URL: &simplexCondition{
			operator: equal,
			value:    "/google-apis",
		},
		QueryParameters: complexConditions{
			complexCondition{
				simplexCondition: simplexCondition{
					operator: equal,
					value:    "MLB",
				},
				field: "site",
			},
		},
		Headers: complexConditions{
			complexCondition{
				simplexCondition: simplexCondition{
					operator: pattern,
					value:    "^Chrome.*",
				},
				field: "User-Agent",
			},
		},
	}
	req := httpRequest{
		URL: "/google-maps-apis",
		QueryParameters: map[string]string{
			"site": "MLB",
		},
		Headers: map[string]string{
			"User-Agent": "Chrome/Google v1.0.9",
		},
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestBodyEquals(t *testing.T) {
	reqMatch := requestMatch{
		Body: &simplexCondition{
			operator: equal,
			value:    `{"status": "PENDING"}`,
		},
	}
	req := httpRequest{
		Body: []byte(`{"status": "PENDING"}`),
	}
	expected := reqMatch.IsExpected(req)
	assert.True(t, expected)
}

func TestBodyNotEquals(t *testing.T) {
	reqMatch := requestMatch{
		Body: &simplexCondition{
			operator: equal,
			value:    `{"status": "PENDING"}`,
		},
	}
	req := httpRequest{
		Body: []byte(`{"status": "CLOSED"}`),
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestBodyContains(t *testing.T) {
	reqMatch := requestMatch{
		Body: &simplexCondition{
			operator: contains,
			value:    "PENDING",
		},
	}
	req := httpRequest{
		Body: []byte(`{"status": "PENDING"}`),
	}
	expected := reqMatch.IsExpected(req)
	assert.True(t, expected)
}

func TestBodyNotContains(t *testing.T) {
	reqMatch := requestMatch{
		Body: &simplexCondition{
			operator: contains,
			value:    "SHIPPED",
		},
	}
	req := httpRequest{
		Body: []byte(`{"status": "PENDING"}`),
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}

func TestBodyPattern(t *testing.T) {
	reqMatch := requestMatch{
		Body: &simplexCondition{
			operator: pattern,
			value:    `^{"[A-Za-z0-9]*":"[A-Za-z0-9]*"}`,
		},
	}
	req := httpRequest{
		Body: []byte(`{"status":"PENDING"}`),
	}
	expected := reqMatch.IsExpected(req)
	assert.True(t, expected)
}

func TestBodyPatternNotMatch(t *testing.T) {
	reqMatch := requestMatch{
		Body: &simplexCondition{
			operator: pattern,
			value:    `^{"[A-Za-z0-9]*":"[A-Za-z0-9]*"}`,
		},
	}
	req := httpRequest{
		Body: []byte(`"status":"PENDING"}`),
	}
	expected := reqMatch.IsExpected(req)
	assert.False(t, expected)
}
