package load_balancer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"github.com/clo-ru/cloapi-go-client/clo/mocks"
	"github.com/clo-ru/cloapi-go-client/clo/request_tools"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sync"
	"testing"
)

func TestBalancerCreateRequest_BuildRequest(t *testing.T) {
	b := BalancerCreateBody{
		Name: "m",
	}
	ID := "id"
	req := BalancerCreateRequest{
		Body:      b,
		ProjectID: ID,
	}
	rawReq, e := req.buildRequest(context.Background(), map[string]interface{}{
		"auth_key": mocks.MockAuthKey,
		"base_url": mocks.MockUrl,
	})
	h := http.Header{}
	h.Add("Authorization", fmt.Sprintf("Bearer %s", mocks.MockAuthKey))
	h.Add("Content-type", "application/json")
	h.Add("X-Add-Some", "SomeHeaderValue")
	rawReq.Header = h
	assert.Nil(t, e)
	bd := new(bytes.Buffer)
	json.NewEncoder(bd).Encode(b)
	expReq, _ := http.NewRequestWithContext(
		context.Background(), http.MethodPost, mocks.MockUrl+fmt.Sprintf(balancerCreateEndpoint, ID), bd,
	)
	expReq.Header = h
	assert.Equal(t, expReq, rawReq)
}

func TestBalancerCreateRequest_Make(t *testing.T) {
	httpCli := mocks.MockClient{}
	cli := clo.ApiClient{
		HttpClient: &httpCli,
		Options: map[string]interface{}{
			"auth_key": "secret",
			"base_url": "https://clo.ru",
		},
	}
	var cases = []struct {
		Name           string
		ShouldFail     bool
		CheckError     bool
		Req            BalancerCreateRequest
		BodyStringFunc func() (string, int)
		Expected       BalancerCreateResponse
	}{
		{
			Name: "Success",
			BodyStringFunc: func() (string, int) {
				return fmt.Sprintf(`{"result":{"id":"id1","algorithm":"algo","addresses":[{"id":"id2","ptr":"ptr2"}],"healthmonitor":{"http_method":"get"}}}`),
					http.StatusOK
			},
			Req: BalancerCreateRequest{
				ProjectID: "id",
			},
			Expected: BalancerCreateResponse{Result: BalancerDetailItem{
				ID:        "id1",
				Algorithm: "algo",
				HealthMonitor: BalancerMonitorDetails{
					HttpMethod: "get",
				},
				Addresses: []BalancerAddress{
					{ID: "id2", Ptr: "ptr2"},
				},
			}},
		},
		{
			Name:       "WrongPtrReturned",
			ShouldFail: true,
			BodyStringFunc: func() (string, int) {
				return fmt.Sprintf(`{"result":{"id":"id1","algorithm":"algo","addresses":[{"id":"id2","ptr":"ptr2"}],"healthmonitor":{"http_method":"get"}}}`),
					http.StatusOK
			},
			Req: BalancerCreateRequest{
				ProjectID: "id",
			},
			Expected: BalancerCreateResponse{Result: BalancerDetailItem{
				ID:        "id1",
				Algorithm: "algo",
				HealthMonitor: BalancerMonitorDetails{
					HttpMethod: "get",
				},
				Addresses: []BalancerAddress{
					{ID: "id2", Ptr: "ptr"},
				},
			}},
		},
		{
			Name:       "Error",
			ShouldFail: true,
			CheckError: true,
			BodyStringFunc: func() (string, int) {
				return "", http.StatusInternalServerError
			},
			Req: BalancerCreateRequest{
				ProjectID: "id",
			},
			Expected: BalancerCreateResponse{},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			mocks.BodyStringFunc = c.BodyStringFunc
			res, e := c.Req.Make(context.Background(), &cli)
			if !c.ShouldFail {
				assert.Nil(t, e)
				assert.Equal(t, c.Expected, res)
			} else {
				if c.CheckError {
					assert.NotNil(t, e)
				} else {
					assert.NotEqual(t, c.Expected, res)
				}
			}
		})
	}
}

func TestBalancerCreateRequest_MakeRetry(t *testing.T) {
	retry := 5
	erCode := http.StatusInternalServerError
	httpCli := mocks.RequestDebugClient{}
	cli := clo.ApiClient{
		HttpClient: &httpCli,
		Options: map[string]interface{}{
			"auth_key": "secret",
			"base_url": "https://clo.ru",
		},
	}
	mocks.BodyStringFunc = func() (string, int) {
		return "", erCode
	}
	grNum := 4
	wg := sync.WaitGroup{}
	for n := 0; n < grNum; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := BalancerCreateRequest{
				Body: BalancerCreateBody{},
			}
			req.WithRetry(retry, 0)
			_, _ = req.Make(context.Background(), &cli)
		}()
	}
	wg.Wait()
	assert.Equal(t, retry*grNum, httpCli.RequestCount)
}

func TestBalancerCreateRequest_CheckPassedBody(t *testing.T) {
	erCode := http.StatusInternalServerError
	httpCli := mocks.RequestDebugClient{}
	cli := clo.ApiClient{
		HttpClient: &httpCli,
		Options: map[string]interface{}{
			"auth_key": "secret",
			"base_url": "https://clo.ru",
		},
	}
	mocks.BodyStringFunc = func() (string, int) {
		return "1", erCode
	}
	req := BalancerCreateRequest{
		Body: BalancerCreateBody{
			Name:               "sname",
			SessionPersistence: true,
			FloatingIP: BalancerBodyAddress{
				ID: "id1",
			},
			HealthMonitor: BalancerBodyMonitor{
				Delay: 3,
			},
			Rules: []BalancerBodyRules{
				{
					PortID: "pid",
				},
			},
		},
	}
	_, _ = req.Make(context.Background(), &cli)
	exp := []byte(`{"name":"sname","algorithm":"","session_persistence":true,"floating_ip":{"id":"id1","ddos_protection":false},"healthmonitor":{"delay":3,"timeout":0,"max_retries":0,"type":""},"rules":[{"port_id":"pid","external_protocol_port":0,"internal_protocol_port":0}]}`)
	exp = append(exp, '\n')

	assert.Equal(t, string(exp), string(httpCli.Body))
}

func TestApiClient_MakeRequestWithError(t *testing.T) {
	mocks.BodyStringFunc = func() (string, int) {
		return `{"code":500,"title":"Internal server error","error":{"message":"try again"}}`,
			http.StatusInternalServerError
	}
	mc := mocks.RequestDebugClient{}
	cli := clo.ApiClient{HttpClient: &mc, Options: map[string]interface{}{"auth_key": "key", "base_url": "1"}}
	req := BalancerCreateRequest{
		ProjectID: "1",
	}
	resp, e := req.Make(context.Background(), &cli)
	assert.Equal(t, BalancerCreateResponse{}, resp)
	assert.Equal(t, request_tools.DefaultError{
		Code:         http.StatusInternalServerError,
		Title:        "Internal server error",
		ErrorMessage: request_tools.ErrorMsg{Message: "try again"},
	}, e)

	resp, e = req.Make(context.Background(), &cli)
	assert.Equal(t, BalancerCreateResponse{}, resp)
	assert.NotEqual(t, request_tools.DefaultError{
		Code:         http.StatusOK,
		Title:        "Internal server error",
		ErrorMessage: request_tools.ErrorMsg{Message: "try again"},
	}, e)
}
