package mock

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Router interface {
	Run(string) error
}

func newRouter(service Service) *router {
	r := &router{
		server:  http.NewServeMux(),
		service: service,
	}
	r.addMappingRoute()
	r.serveMockRoute()
	return r
}

type router struct {
	server  *http.ServeMux
	service Service
}

func (r *router) Run(address string) error {
	return http.ListenAndServe(address, r.server)
}

func (r *router) addMappingRoute() {
	r.server.HandleFunc("/mock/mapping", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var dto mockDTO
		err := decodeAsJson(request.Body, &dto)
		if err != nil {
			writeErrorAsJson(err, writer)
			return
		}
		resp, err := r.service.Add(dto)
		if err != nil {
			writeErrorAsJson(err, writer)
			return
		}
		writeAsJson(writer, resp, http.StatusOK)
		return
	})
}

func (r *router) serveMockRoute() {
	r.server.HandleFunc("/", func(writer http.ResponseWriter, httpRequest *http.Request) {
		request := buildRequest(httpRequest)
		resp, err := r.service.Match(request)
		if err != nil {
			writeErrorAsJson(err, writer)
			return
		}
		writeHttpResponse(writer, resp)
		return
	})
}

func writeErrorAsJson(err error, writer http.ResponseWriter) {
	switch t := err.(type) {
	case Error:
		{
			writeAsJson(writer, err, getHttpStatusCodeBy(t.Code))
			break
		}
	default:
		{
			writeAsJson(writer, err, http.StatusInternalServerError)
			break
		}
	}
}

func writeAsJson(writer http.ResponseWriter, resp any, status int) {
	data, _ := json.Marshal(resp)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	_, err := writer.Write(data)
	if err != nil {
		LogError("error writing json http response, error: %v", err)
	}
}

func writeHttpResponse(writer http.ResponseWriter, response *httpResponse) {
	for key, value := range response.Headers {
		writer.Header().Add(key, value)
	}
	writer.WriteHeader(response.Status)
	_, err := writer.Write(response.Body)
	if err != nil {
		LogError("error writing http response, error: %v", err)
	}
}

func decodeAsJson(body io.ReadCloser, destination interface{}) error {
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(buffer.Bytes(), destination)
}

func buildRequest(request *http.Request) httpRequest {
	queryParams := request.URL.Query()
	header := request.Header
	buf := new(bytes.Buffer)
	if request.Body != nil {
		buf.ReadFrom(request.Body)
	}
	return httpRequest{
		URL:             request.URL.Path,
		Method:          request.Method,
		QueryParameters: flatValues(queryParams),
		Headers:         flatValues(header),
		Body:            buf.Bytes(),
	}
}

func flatValues(data map[string][]string) map[string]string {
	flatMap := map[string]string{}
	for key, value := range data {
		flatMap[key] = strings.Join(value, ",")
	}
	return flatMap
}

func getHttpStatusCodeBy(domainCode string) int {
	switch domainCode {
	case "invalid_request":
		return http.StatusBadRequest
	case "mock_not_found":
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
