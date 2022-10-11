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

func TestFipPtrChangeRequest_BuildRequest(t *testing.T) {
	b := FipPtrChangeBody{
		Value: "id",
	}
	fID := "id"
	req := FipPtrChangeRequest{
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
		context.Background(), http.MethodPut, mocks.MockUrl+fmt.Sprintf(fipPtrChangeEndpoint, fID), bd,
	)
	expReq.Header = h
	assert.Equal(t, expReq, rawReq)
}

func TestFipPtrChangeRequest_Make(t *testing.T) {
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
		Req            FipPtrChangeRequest
		BodyStringFunc func() (string, int)
	}{
		{
			Name: "Success",
			BodyStringFunc: func() (string, int) {
				return "1",
					http.StatusOK
			},
			Req: FipPtrChangeRequest{
				FloatingID: "id",
			},
		},
		{
			Name:       "Error",
			ShouldFail: true,
			CheckError: true,
			BodyStringFunc: func() (string, int) {
				return "", http.StatusInternalServerError
			},
			Req: FipPtrChangeRequest{
				FloatingID: "id",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			mocks.BodyStringFunc = c.BodyStringFunc
			e := c.Req.Make(context.Background(), &cli)
			if !c.ShouldFail {
				assert.Nil(t, e)
			} else {
				assert.NotNil(t, e)
			}
		})
	}
}

func TestFipPtrChangeRequest_MakeRetry(t *testing.T) {
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
			req := FipPtrChangeRequest{
				Body: FipPtrChangeBody{
					Value: "new_value",
				},
			}
			req.WithRetry(retry, 0)
			_ = req.Make(context.Background(), &cli)
		}()
	}
	wg.Wait()
	assert.Equal(t, retry*grNum, httpCli.RequestCount)
}

func TestFipPtrChangeRequest_CheckPassedBody(t *testing.T) {
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
	req := FipPtrChangeRequest{
		Body: FipPtrChangeBody{
			Value: "new_ptr",
		},
	}
	_ = req.Make(context.Background(), &cli)
	exp := []byte(`{"value":"new_ptr"}`)
	exp = append(exp, '\n')

	assert.Equal(t, exp, httpCli.Body)
}
