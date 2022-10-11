package storage

import (
	"github.com/clo-ru/cloapi-go-client/clo"
	"github.com/clo-ru/cloapi-go-client/clo/mocks"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sync"
	"testing"
)

func TestS3UserUnsuspendRequest_BuildRequest(t *testing.T) {
	ID := "id"
	req := S3UserUnsuspendRequest{
		UserID: ID,
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
		context.Background(), http.MethodPost, mocks.MockUrl+fmt.Sprintf(s3UserUnsuspendEndpoint, ID), nil,
	)
	expReq.Header = h
	assert.Equal(t, expReq, rawReq)
}

func TestS3UserUnsuspendRequest_Make(t *testing.T) {
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
		Req            S3UserUnsuspendRequest
		BodyStringFunc func() (string, int)
		Expected       S3UserUnsuspendResponse
	}{
		{
			Name: "Success",
			BodyStringFunc: func() (string, int) {
				return fmt.Sprintf(`{"result":{"id":"user_id","max_buckets":2,"quotas":[{"max_size":1,"max_objects":2,"type":"user"}]}}`),
					http.StatusOK
			},
			Req: S3UserUnsuspendRequest{},
			Expected: S3UserUnsuspendResponse{Result: ResponseItem{
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
			Req: S3UserUnsuspendRequest{},
			Expected: S3UserUnsuspendResponse{Result: ResponseItem{
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
			Req:      S3UserUnsuspendRequest{},
			Expected: S3UserUnsuspendResponse{},
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

func TestS3UserUnsuspendRequest_MakeRetry(t *testing.T) {
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
			req := S3UserUnsuspendRequest{}
			req.WithRetry(retry, 0)
			_, _ = req.Make(context.Background(), &cli)
		}()
	}
	wg.Wait()
	assert.Equal(t, retry*grNum, httpCli.RequestCount)
}
