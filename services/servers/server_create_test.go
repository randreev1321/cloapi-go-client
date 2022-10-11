package servers

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

func TestServerCreateRequest_BuildRequest(t *testing.T) {
	b := ServerCreateBody{
		Name: "m",
	}
	ID := "id"
	req := ServerCreateRequest{
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
		context.Background(), http.MethodPost, mocks.MockUrl+fmt.Sprintf(serverCreateEndpoint, ID), bd,
	)
	expReq.Header = h
	assert.Equal(t, expReq, rawReq)
}

func TestServerCreateRequest_Make(t *testing.T) {
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
		Req            ServerCreateRequest
		BodyStringFunc func() (string, int)
		Expected       ServerCreateResponse
	}{
		{
			Name: "Success",
			BodyStringFunc: func() (string, int) {
				return fmt.Sprintf(`{"result":{"id":"sid","flavor":{"ram":2,"vcpus":3}}}`),
					http.StatusOK
			},
			Req: ServerCreateRequest{
				ProjectID: "id",
			},
			Expected: ServerCreateResponse{Result: ServerCreateItem{
				ID: "sid",
				Flavor: ServerFlavorData{
					Ram:   2,
					Vcpus: 3,
				},
			}},
		},
		{
			Name:       "WrongRamReturned",
			ShouldFail: true,
			BodyStringFunc: func() (string, int) {
				return fmt.Sprintf(`{"result":{"id":"sid","flavor":{"ram":1,"vcpus":3}}}`),
					http.StatusOK
			},
			Req: ServerCreateRequest{
				ProjectID: "id",
			},
			Expected: ServerCreateResponse{Result: ServerCreateItem{
				ID: "sid",
				Flavor: ServerFlavorData{
					Ram:   2,
					Vcpus: 3,
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
			Req: ServerCreateRequest{
				ProjectID: "id",
			},
			Expected: ServerCreateResponse{},
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

func TestServerCreateRequest_MakeRetry(t *testing.T) {
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
			req := ServerCreateRequest{
				Body: ServerCreateBody{},
			}
			req.WithRetry(retry, 0)
			_, _ = req.Make(context.Background(), &cli)
		}()
	}
	wg.Wait()
	assert.Equal(t, retry*grNum, httpCli.RequestCount)
}

func TestServerCreateRequest_CheckPassedBody(t *testing.T) {
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
	req := ServerCreateRequest{
		Body: ServerCreateBody{
			Name:     "sname",
			Image:    "image",
			Recipe:   "recipe",
			Keypairs: []string{"pubkey"},
			Flavor:   ServerFlavorBody{Ram: 3, Vcpus: 2},
			Storages: []ServerStorageBody{
				{
					Size:        10,
					Bootable:    true,
					StorageType: "volume",
				},
			},
			Licenses: []ServerLicenseBody{
				{
					Value: 1,
					Name:  "ispmgr",
				},
			},
			Addresses: []ServerAddressesBody{
				{
					External:       true,
					DdosProtection: true,
					FloatingIpID:   "123",
					Version:        4,
				},
			},
		},
	}
	_, _ = req.Make(context.Background(), &cli)
	exp := []byte(`{"name":"sname","image":"image","recipe":"recipe","keypairs":["pubkey"],"flavor":{"ram":3,"vcpus":2},"storages":[{"size":10,"bootable":true,"storage_type":"volume"}],"licenses":[{"value":1,"name":"ispmgr"}],"addresses":[{"external":true,"ddos_protection":true,"floatingip_id":"123","version":4}]}`)
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
	req := ServerCreateRequest{
		ProjectID: "1",
	}
	resp, e := req.Make(context.Background(), &cli)
	assert.Equal(t, ServerCreateResponse{}, resp)
	assert.Equal(t, request_tools.DefaultError{
		Code:         http.StatusInternalServerError,
		Title:        "Internal server error",
		ErrorMessage: request_tools.ErrorMsg{Message: "try again"},
	}, e)

	resp, e = req.Make(context.Background(), &cli)
	assert.Equal(t, ServerCreateResponse{}, resp)
	assert.NotEqual(t, request_tools.DefaultError{
		Code:         http.StatusOK,
		Title:        "Internal server error",
		ErrorMessage: request_tools.ErrorMsg{Message: "try again"},
	}, e)
}
