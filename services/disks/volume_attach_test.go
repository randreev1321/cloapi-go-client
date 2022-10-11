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

func TestVolumeAttachRequest_BuildRequest(t *testing.T) {
	b := VolumeAttachBody{
		MountPath: "m",
		ServerID:  "id",
	}
	dID := "volume_id"
	req := VolumeAttachRequest{
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
		context.Background(), http.MethodPost, mocks.MockUrl+fmt.Sprintf(volumeAttachEndpoint, dID), bd,
	)
	expReq.Header = h
	assert.Equal(t, expReq, rawReq)
}

func TestVolumeAttachRequest_Make(t *testing.T) {
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
		Req            VolumeAttachRequest
		BodyStringFunc func() (string, int)
		Expected       VolumeAttachResponse
	}{
		{
			Name: "Success",
			BodyStringFunc: func() (string, int) {
				return fmt.Sprintf(`{"result":{"device":"sda", "volume":{"id": "disk_id"}, "server":{"id":"server"}}}`),
					http.StatusOK
			},
			Req: VolumeAttachRequest{
				VolumeID: "disk_id",
			},
			Expected: VolumeAttachResponse{Result: VolumeAttachItem{
				Device: "sda",
				Volume: VolumeAttachDetail{
					ID: "disk_id",
				},
				Server: AttachedToServer{
					ID: "server",
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
			Req: VolumeAttachRequest{
				VolumeID: "disk_id",
			},
			Expected: VolumeAttachResponse{Result: VolumeAttachItem{}},
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

func TestVolumeAttachRequest_MakeRetry(t *testing.T) {
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
			req := VolumeAttachRequest{
				Body: VolumeAttachBody{
					MountPath: "/var/opt",
					ServerID:  "server_id",
				},
			}
			req.WithRetry(retry, 0)
			_, _ = req.Make(context.Background(), &cli)
		}()
	}
	wg.Wait()
	assert.Equal(t, retry*grNum, httpCli.RequestCount)
}

func TestVolumeAttachRequest_CheckPassedBody(t *testing.T) {
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
	req := VolumeAttachRequest{
		Body: VolumeAttachBody{
			MountPath: "/var/opt",
			ServerID:  "server_id",
		},
	}
	_, _ = req.Make(context.Background(), &cli)
	exp := []byte(`{"mount_path":"/var/opt","server_id":"server_id"}`)
	exp = append(exp, '\n')

	assert.Equal(t, exp, httpCli.Body)
}
