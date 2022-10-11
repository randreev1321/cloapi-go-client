package disks

import (
	"bytes"
	"github.com/clo-ru/cloapi-go-client/clo"
	"github.com/clo-ru/cloapi-go-client/clo/mocks"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sync"
	"testing"
)

func TestVolumeCreateRequest_BuildRequest(t *testing.T) {
	b := VolumeCreateBody{
		Name: "m",
	}
	dID := "volume_id"
	req := VolumeCreateRequest{
		Body:      b,
		ProjectID: dID,
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
		context.Background(), http.MethodPost, mocks.MockUrl+fmt.Sprintf(volumeCreateEndpoint, dID), bd,
	)
	expReq.Header = h
	assert.Equal(t, expReq, rawReq)
}

func TestVolumeCreateRequest_Make(t *testing.T) {
	httpCli := mocks.RequestDebugClient{}
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
		Req            VolumeCreateRequest
		BodyStringFunc func() (string, int)
		Expected       VolumeCreateResponse
	}{
		{
			Name: "Success",
			BodyStringFunc: func() (string, int) {
				return fmt.Sprintf(`{"result":{"id":"disk_id","name":"some_name","device":"sda","undetachable":true,"attached_to_server":{"id":"server"}}}`),
					http.StatusOK
			},
			Req: VolumeCreateRequest{
				ProjectID: "project_id",
			},
			Expected: VolumeCreateResponse{
				Result: VolumeDetail{
					Device:       "sda",
					Undetachable: true,
					ResponseItem: ResponseItem{
						ID:               "disk_id",
						Name:             "some_name",
						AttachedToServer: AttachedToServer{ID: "server"},
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
			Req: VolumeCreateRequest{
				ProjectID: "project_id",
			},
			Expected: VolumeCreateResponse{},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			mocks.BodyStringFunc = c.BodyStringFunc
			res, e := c.Req.Make(context.Background(), &cli)
			assert.True(t, mocks.CheckHeaders(httpCli.Headers))
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

func TestVolumeCreateRequest_MakeRetry(t *testing.T) {
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
			req := VolumeCreateRequest{
				Body: VolumeCreateBody{},
			}
			req.WithRetry(retry, 0)
			_, _ = req.Make(context.Background(), &cli)
		}()
	}
	wg.Wait()
	assert.Equal(t, retry*grNum, httpCli.RequestCount)
}

func TestVolumeCreateRequest_CheckPassedBody(t *testing.T) {
	erCode := http.StatusOK
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
	req := VolumeCreateRequest{
		Body: VolumeCreateBody{
			Name: "sda",
			Size: 11,
		},
	}
	_, _ = req.Make(context.Background(), &cli)
	exp := []byte(`{"name":"sda","size":11,"autorename":false}`)
	exp = append(exp, '\n')

	assert.Equal(t, exp, httpCli.Body)
}
