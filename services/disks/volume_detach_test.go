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

func TestVolumeDetachRequest_BuildRequest(t *testing.T) {
	b := VolumeDetachBody{
		Force: false,
	}
	dID := "volume_id"
	req := VolumeDetachRequest{
		Body:     b,
		VolumeID: dID,
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
		context.Background(), http.MethodPost, mocks.MockUrl+fmt.Sprintf(volumeDetachEndpoint, dID), bd,
	)
	expReq.Header = h
	assert.Equal(t, expReq, rawReq)
}

func TestVolumeDetachRequest_Make(t *testing.T) {
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
		Req            VolumeDetachRequest
		BodyStringFunc func() (string, int)
	}{
		{
			Name: "Success",
			BodyStringFunc: func() (string, int) {
				return "1",
					http.StatusOK
			},
			Req: VolumeDetachRequest{
				VolumeID: "disk_id",
			},
		},
		{
			Name:       "Error",
			ShouldFail: true,
			BodyStringFunc: func() (string, int) {
				return "", http.StatusInternalServerError
			},
			Req: VolumeDetachRequest{
				VolumeID: "disk_id",
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

func TestVolumeDetachRequest_MakeRetry(t *testing.T) {
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
			req := VolumeDetachRequest{
				Body: VolumeDetachBody{
					Force: true,
				},
			}
			req.WithRetry(retry, 0)
			_ = req.Make(context.Background(), &cli)
		}()
	}
	wg.Wait()
	assert.Equal(t, retry*grNum, httpCli.RequestCount)
}

func TestVolumeDetachRequest_CheckPassedBody(t *testing.T) {
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
	req := VolumeDetachRequest{
		Body: VolumeDetachBody{
			Force: true,
		},
	}
	_ = req.Make(context.Background(), &cli)

	exp := []byte(`{"force":true}`)
	exp = append(exp, '\n')
	assert.Equal(t, exp, httpCli.Body)

	exp = []byte(`{"force":false}`)
	exp = append(exp, '\n')
	assert.NotEqual(t, exp, httpCli.Body)

}
