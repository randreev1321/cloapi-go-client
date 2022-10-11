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

func TestFipAttachRequest_BuildRequest(t *testing.T) {
	b := FipAttachBody{
		ServerID: "id",
	}
	fID := "id"
	req := FipServerAttachRequest{
		Body:       b,
		FloatingID: fID,
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
		context.Background(), http.MethodPost, mocks.MockUrl+fmt.Sprintf(fipServerAttachEndpoint, fID), bd,
	)
	expReq.Header = h
	assert.Equal(t, expReq, rawReq)
}

func TestFipAttachRequest_Make(t *testing.T) {
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
		Req            FipServerAttachRequest
		BodyStringFunc func() (string, int)
		Expected       FipServerAttachResponse
	}{
		{
			Name: "Success",
			BodyStringFunc: func() (string, int) {
				return fmt.Sprintf(`{"result":{"id":"fipid","status":"ACTIVE","attached_to_server":{"id":"server_id"}}}`),
					http.StatusOK
			},
			Req: FipServerAttachRequest{
				FloatingID: "id",
			},
			Expected: FipServerAttachResponse{Result: FipDetail{
				ID:     "fipid",
				Status: "ACTIVE",
				AttachedToServer: AttachedToEntityDetails{
					ID: "server_id",
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
			Req: FipServerAttachRequest{
				FloatingID: "id",
			},
			Expected: FipServerAttachResponse{Result: FipDetail{}},
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

func TestFipAttachRequest_MakeRetry(t *testing.T) {
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
			req := FipServerAttachRequest{
				Body: FipAttachBody{
					ServerID: "server_id",
				},
			}
			req.WithRetry(retry, 0)
			_, _ = req.Make(context.Background(), &cli)
		}()
	}
	wg.Wait()
	assert.Equal(t, retry*grNum, httpCli.RequestCount)
}

func TestFipAttachRequest_CheckPassedBody(t *testing.T) {
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
	req := FipServerAttachRequest{
		Body: FipAttachBody{
			ServerID: "server_id",
		},
	}
	_, _ = req.Make(context.Background(), &cli)
	exp := []byte(`{"id":"server_id"}`)
	exp = append(exp, '\n')

	assert.Equal(t, exp, httpCli.Body)
}
