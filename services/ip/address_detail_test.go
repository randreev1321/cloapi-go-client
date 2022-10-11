package ip

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

func TestAddressDetailRequest_BuildRequest(t *testing.T) {
	ID := "id"
	req := AddressDetailRequest{
		AddressID: ID,
	}
	h := http.Header{}
	h.Add("Authorization", fmt.Sprintf("Bearer %s", mocks.MockAuthKey))
	h.Add("Content-type", "application/json")
	h.Add("X-Add-Some", "SomeHeaderValue")
	rawReq, e := req.buildRequest(context.Background(), map[string]interface{}{
		"auth_key": mocks.MockAuthKey,
		"base_url": mocks.MockUrl,
	})
	rawReq.Header = h
	assert.Nil(t, e)
	expReq, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet, mocks.MockUrl+fmt.Sprintf(addressDetailEndpoint, ID), nil,
	)
	expReq.Header = h
	assert.Equal(t, expReq, rawReq)
}

func TestAddressDetailRequest_Make(t *testing.T) {
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
		Req            AddressDetailRequest
		BodyStringFunc func() (string, int)
		Expected       AddressDetailResponse
	}{
		{
			Name: "Success",
			BodyStringFunc: func() (string, int) {
				return fmt.Sprintf(`{"result":{"id":"fipid","status":"ACTIVE","address":"192.168.1.1"}}`),
					http.StatusOK
			},
			Req: AddressDetailRequest{
				AddressID: "id",
			},
			Expected: AddressDetailResponse{Result: AddressDetail{
				ID:      "fipid",
				Status:  "ACTIVE",
				Address: "192.168.1.1",
			}},
		},
		{
			Name:       "Error",
			ShouldFail: true,
			CheckError: true,
			BodyStringFunc: func() (string, int) {
				return "", http.StatusInternalServerError
			},
			Req: AddressDetailRequest{
				AddressID: "id",
			},
			Expected: AddressDetailResponse{Result: AddressDetail{}},
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

func TestAddressDetailRequest_MakeRetry(t *testing.T) {
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
			req := AddressDetailRequest{}
			req.WithRetry(retry, 0)
			_, _ = req.Make(context.Background(), &cli)
		}()
	}
	wg.Wait()
	assert.Equal(t, retry*grNum, httpCli.RequestCount)
}
