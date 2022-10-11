package storage

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

func TestS3UserPatchRequest_BuildRequest(t *testing.T) {
	b := S3UserPatchBody{
		Name: "m",
	}
	ID := "id"
	req := S3UserPatchRequest{
		UserID: ID,
		Body:   b,
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
		context.Background(), http.MethodPatch, mocks.MockUrl+fmt.Sprintf(s3UserPatchEndpoint, ID), bd,
	)
	expReq.Header = h
	assert.Equal(t, expReq, rawReq)
}

func TestS3UserPatchRequest_Make(t *testing.T) {
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
		Req            S3UserPatchRequest
		BodyStringFunc func() (string, int)
		Expected       S3UserPatchResponse
	}{
		{
			Name: "Success",
			BodyStringFunc: func() (string, int) {
				return fmt.Sprintf(`{"result":{"id":"user_id","max_buckets":2,"quotas":[{"max_size":1,"max_objects":2,"type":"user"}]}}`),
					http.StatusOK
			},
			Req: S3UserPatchRequest{},
			Expected: S3UserPatchResponse{Result: ResponseItem{
				ID:         "user_id",
				MaxBuckets: 2,
				Quotas: []QuotaInfo{
					{
						MaxSize:    1,
						MaxObjects: 2,
						Type:       "user",
					},
				},
			}},
		},
		{
			Name:       "WrongMaxBucketsReturned",
			ShouldFail: true,
			BodyStringFunc: func() (string, int) {
				return fmt.Sprintf(`{"result":{"id":"user_id","max_buckets":1,"quotas":[{"max_size":1,"max_objects":2,"type":"user"}]}}`),
					http.StatusOK
			},
			Req: S3UserPatchRequest{},
			Expected: S3UserPatchResponse{Result: ResponseItem{
				ID:         "user_id",
				MaxBuckets: 2,
				Quotas: []QuotaInfo{
					{
						MaxSize:    1,
						MaxObjects: 2,
						Type:       "user",
					},
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
			Req:      S3UserPatchRequest{},
			Expected: S3UserPatchResponse{},
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

func TestS3UserPatchRequest_MakeRetry(t *testing.T) {
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
			req := S3UserPatchRequest{
				Body: S3UserPatchBody{
					Name: "name",
				},
			}
			req.WithRetry(retry, 0)
			_, _ = req.Make(context.Background(), &cli)
		}()
	}
	wg.Wait()
	assert.Equal(t, retry*grNum, httpCli.RequestCount)
}

func TestS3UserPatchRequest_CheckPassedBody(t *testing.T) {
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
	req := S3UserPatchRequest{
		Body: S3UserPatchBody{
			Name: "new_name",
		},
	}
	_, _ = req.Make(context.Background(), &cli)
	exp := []byte(`{"name":"new_name"}`)
	exp = append(exp, '\n')
	assert.Equal(t, exp, httpCli.Body)

	exp = []byte(`{"name":"ew_name"}`)
	exp = append(exp, '\n')
	assert.NotEqual(t, exp, httpCli.Body)
}
