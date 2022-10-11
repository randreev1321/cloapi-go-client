package ip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"github.com/clo-ru/cloapi-go-client/clo/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sync"
	"testing"
)

func TestAddressAttachRequest_BuildRequest(t *testing.T) {
	b := AddressAttachBody{
		ID:     "id",
		Entity: "server",
	}
	fID := "id"
	req := AddressAttachRequest{
		Body:      b,
		AddressID: fID,
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
	bd := new(bytes.Buffer)
	json.NewEncoder(bd).Encode(b)
	expReq, _ := http.NewRequestWithContext(
		context.Background(), http.MethodPost, mocks.MockUrl+fmt.Sprintf(addressAttachEndpoint, fID), bd,
	)
	expReq.Header = h
	assert.Equal(t, expReq, rawReq)
}

func TestAddressAttachRequest_Make(t *testing.T) {
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
		Req            AddressAttachRequest
		BodyStringFunc func() (string, int)
		Expected       AddressAttachResponse
	}{
		{
			Name: "Success",
			BodyStringFunc: func() (string, int) {
				return fmt.Sprintf(`{"result":{"id":"fipid","status":"ACTIVE","attached_to":{"id":"server_id","entity":"load_balancer"}}}`),
					http.StatusOK
			},
			Req: AddressAttachRequest{
				AddressID: "id",
			},
			Expected: AddressAttachResponse{Result: AddressDetail{
				ID:     "fipid",
				Status: "ACTIVE",
				AttachedTo: AttachedToDetails{
					Entity: "load_balancer",
					ID:     "server_id",
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
			Req: AddressAttachRequest{
				AddressID: "id",
			},
			Expected: AddressAttachResponse{Result: AddressDetail{}},
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

func TestAddressAttachRequest_MakeRetry(t *testing.T) {
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
			req := AddressAttachRequest{
				Body: AddressAttachBody{
					Entity: "server",
					ID:     "server_id",
				},
			}
			req.WithRetry(retry, 0)
			_, _ = req.Make(context.Background(), &cli)
		}()
	}
	wg.Wait()
	assert.Equal(t, retry*grNum, httpCli.RequestCount)
}

func TestAddressAttachRequest_CheckPassedBody(t *testing.T) {
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
	req := AddressAttachRequest{
		Body: AddressAttachBody{
			Entity: "server",
			ID:     "server_id",
		},
	}
	_, _ = req.Make(context.Background(), &cli)
	exp := []byte(`{"id":"server_id","entity":"server"}`)
	exp = append(exp, '\n')

	assert.Equal(t, exp, httpCli.Body)
}
