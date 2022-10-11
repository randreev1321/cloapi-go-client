package clo

import (
	"fmt"
	"net/http"
	"time"
)

type Request struct {
	retry        int
	retryTimeout time.Duration
	log          Logger
	headers      http.Header
	queryParams  map[string][]string
}

func (r *Request) WithLog(l Logger) {
	r.log = l
}

func (r *Request) WithRetry(retry int, timeout time.Duration) {
	r.retry = retry
	r.retryTimeout = timeout
}

func (r *Request) WithQueryParams(param map[string][]string) {
	if r.queryParams == nil {
		r.queryParams = param
	} else {
		for k, v := range param {
			if _, ok := r.queryParams[k]; !ok {
				r.queryParams[k] = []string{}
			}
			r.queryParams[k] = append(r.queryParams[k], v...)
		}
	}
}

func (r *Request) WithHeaders(headers map[string][]string) {
	if r.headers == nil {
		r.headers = map[string][]string{}
	}
	r.headers = headers
}

func (r *Request) MakeRequest(req *http.Request, cli *ApiClient) (*http.Response, error) {
	r.addUrlParams(req)
	r.addHeadersToReq(req)
	if r.retry < 0 {
		return nil, fmt.Errorf("retry number should be positive")
	}
	var (
		rawResp      *http.Response
		requestError error
	)
	for r.retry >= 0 {
		rawResp, requestError = cli.MakeRequest(req)
		if requestError == nil {
			r.retry = -1
			break
		}
		r.retry -= 1
		if r.log != nil {
			r.log.Errorf("%T error %s, the request will be retried %d more times\n",
				*r, requestError.Error(), r.retry)
		}
		if r.retry == 0 {
			break
		}
		time.Sleep(r.retryTimeout)
	}
	if requestError != nil {
		return nil, requestError
	}
	return rawResp, nil
}

func (r *Request) addUrlParams(req *http.Request) {
	rawQuery := req.URL.Query()
	for pn, pv := range r.queryParams {
		for _, v := range pv {
			rawQuery.Add(pn, v)
		}
	}
	req.URL.RawQuery = rawQuery.Encode()
}

func (r *Request) addHeadersToReq(req *http.Request) {
	req.Header = r.headers
	if req.Header == nil {
		req.Header = http.Header{}
	}
	req.Header.Add("Content-type", "application/json")
}
