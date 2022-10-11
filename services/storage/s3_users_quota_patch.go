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
	s3UserQuotaPatchEndpoint = "/v1/s3_users/%s/quotas"
)

type S3UserQuotaPatchRequest struct {
	clo.Request
	UserID string
	Body   S3UserQuotaPatchBody
}

type S3UserQuotaPatchBody struct {
	MaxBuckets  int              `json:"max_buckets"`
	UserQuota   PatchQuotaParams `json:"user_quota"`
	BucketQuota PatchQuotaParams `json:"bucket_quota"`
}

type PatchQuotaParams struct {
	MaxObjects int `json:"max_objects"`
	MaxSize    int `json:"max_size"`
}

func (r *S3UserQuotaPatchRequest) Make(ctx context.Context, cli *clo.ApiClient) (S3UserQuotaPatchResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return S3UserQuotaPatchResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return S3UserQuotaPatchResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp S3UserQuotaPatchResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return S3UserQuotaPatchResponse{}, e
	}
	return resp, nil
}

func (r *S3UserQuotaPatchRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(s3UserQuotaPatchEndpoint, r.UserID)
	b := new(bytes.Buffer)
	if e := json.NewEncoder(b).Encode(r.Body); e != nil {
		return nil, fmt.Errorf("can't encode body parameters, %s", e.Error())
	}
	rawReq, e := http.NewRequestWithContext(
		ctx, http.MethodPatch, baseUrl, b,
	)
	if e != nil {
		return nil, e
	}
	h := http.Header{}
	h.Add("Authorization", fmt.Sprintf("Bearer %s", authKey))
	r.WithHeaders(h)
	return rawReq, nil
}
