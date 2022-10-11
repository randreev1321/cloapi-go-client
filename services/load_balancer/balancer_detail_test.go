package load_balancer

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"github.com/clo-ru/cloapi-go-client/clo/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sync"
	"testing"
)

func TestBalancerDetailRequest_BuildRequest(t *testing.T) {
	ID := "id"
	req := BalancerDetailRequest{
		ServerID: ID,
	}
	h := http.Header{}
	h.Add("X-Add-Some", "SomeHeaderValue")
	h.Add("Authorization", fmt.Sprintf("Bearer %s", mocks.MockAuthKey))
	h.Add("Content-type", "application/json")
	rawReq, e := req.buildRequest(context.Background(), map[string]interface{}{
		"auth_key": mocks.MockAuthKey,
		"base_url": mocks.MockUrl,
	})
	assert.Nil(t, e)
	rawReq.Header = h
	expReq, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet, mocks.MockUrl+fmt.Sprintf(balancerDetailEndpoint, ID), nil,
	)
	expReq.Header = h
	assert.Equal(t, expReq, rawReq)
}

func TestBalancerDetailRequest_Make(t *testing.T) {
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
		Req            BalancerDetailRequest
		BodyStringFunc func() (string, int)
		Expected       BalancerDetailResponse
	}{
		{
			Name: "Success",
			BodyStringFunc: func() (string, int) {
				return fmt.Sprintf(`{"result":{"id":"id1","algorithm":"algo","addresses":[{"id":"id2","ptr":"ptr2"}],"healthmonitor":{"http_method":"get"}}}`),
					http.StatusOK
			},
			Req: BalancerDetailRequest{
				ServerID: "id",
			},
			Expected: BalancerDetailResponse{Result: BalancerDetailItem{
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
			Req: BalancerDetailRequest{
				ServerID: "id",
			},
			Expected: BalancerDetailResponse{Result: BalancerDetailItem{
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
			Req: BalancerDetailRequest{
				ServerID: "id",
			},
			Expected: BalancerDetailResponse{},
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

func TestBalancerDetailRequest_MakeRetry(t *testing.T) {
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
			req := BalancerDetailRequest{}
			req.WithRetry(retry, 0)
			_, _ = req.Make(context.Background(), &cli)
		}()
	}
	wg.Wait()
	assert.Equal(t, retry*grNum, httpCli.RequestCount)
}
