package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	s3UserCreateEndpoint = "/v1/projects/%s/s3_users"
)

type S3UserCreateRequest struct {
	clo.Request
	ProjectID string
	Body      S3UserCreateBody
}

type S3UserCreateBody struct {
	Name          string            `json:"name"`
	CanonicalName string            `json:"canonical_name"`
	DefaultBucket bool              `json:"default_bucket"`
	MaxBuckets    int               `json:"max_buckets"`
	BucketQuota   CreateQuotaParams `json:"bucket_quota,omitempty"`
	UserQuota     CreateQuotaParams `json:"user_quota"`
}

type CreateQuotaParams struct {
	MaxObjects int `json:"max_objects,omitempty"`
	MaxSize    int `json:"max_size"`
}

func (r *S3UserCreateRequest) Make(ctx context.Context, cli *clo.ApiClient) (S3UserCreateResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return S3UserCreateResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return S3UserCreateResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp S3UserCreateResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return S3UserCreateResponse{}, e
	}
	return resp, nil
}

func (r *S3UserCreateRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(s3UserCreateEndpoint, r.ProjectID)
	b := new(bytes.Buffer)
	if e := json.NewEncoder(b).Encode(r.Body); e != nil {
		return nil, fmt.Errorf("can't encode body parameters, %s", e.Error())
	}
	rawReq, e := http.NewRequestWithContext(
		ctx, http.MethodPost, baseUrl, b,
	)
	if e != nil {
		return nil, e
	}
	h := http.Header{}
	h.Add("Authorization", fmt.Sprintf("Bearer %s", authKey))
	r.WithHeaders(h)
	return rawReq, nil
}
