package mocks

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	MockUrl     = "https://clo.ru"
	MockAuthKey = "supersecret"
)

var BodyStringFunc func() (string, int)

type MockClient struct {
	RequestCount int
}

func (mc *MockClient) Do(req *http.Request) (*http.Response, error) {
	mc.RequestCount++
	var bodyString string
	var statusCode int
	bodyString, statusCode = BodyStringFunc()
	if len(bodyString) == 0 {
		return nil, fmt.Errorf(http.StatusText(statusCode))
	}
	body := ioutil.NopCloser(bytes.NewReader([]byte(bodyString)))
	return &http.Response{StatusCode: statusCode, Body: body}, nil
}

//RequestDebugClient is useful when you want to discover the passed URL and parameters
type RequestDebugClient struct {
	RequestCount int
	URL          url.URL
	Headers      http.Header
	Body         []byte
}

func (rdc *RequestDebugClient) Do(req *http.Request) (*http.Response, error) {
	rdc.RequestCount++
	rdc.URL = *req.URL
	rdc.Headers = req.Header
	var e error
	if req.Body != nil {
		if rdc.Body, e = io.ReadAll(req.Body); e != nil {
			return nil, e
		}
		if e = req.Body.Close(); e != nil {
			return nil, e
		}
	}
	var bodyString string
	var statusCode int
	bodyString, statusCode = BodyStringFunc()
	if len(bodyString) == 0 {
		return nil, fmt.Errorf(http.StatusText(statusCode))
	}
	body := ioutil.NopCloser(bytes.NewReader([]byte(bodyString)))
	return &http.Response{StatusCode: statusCode, Body: body}, nil
}

func CheckHeaders(headers map[string][]string) bool {
	if _, ok := headers["Authorization"]; !ok {
		return false
	}
	if _, ok := headers["Content-Type"]; !ok {
		return false
	}
	return true
}
