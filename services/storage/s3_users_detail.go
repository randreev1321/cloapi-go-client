package storage

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	s3UserDetailEndpoint = "/v1/s3_users/%s"
)

type S3UserDetailRequest struct {
	clo.Request
	UserID string
}

func (r *S3UserDetailRequest) Make(ctx context.Context, cli *clo.ApiClient) (S3UserDetailResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return S3UserDetailResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return S3UserDetailResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp S3UserDetailResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return S3UserDetailResponse{}, e
	}
	return resp, nil
}

func (r *S3UserDetailRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(s3UserDetailEndpoint, r.UserID)
	rawReq, e := http.NewRequestWithContext(
		ctx, http.MethodGet, baseUrl, nil,
	)
	if e != nil {
		return nil, e
	}
	h := http.Header{}
	h.Add("Authorization", fmt.Sprintf("Bearer %s", authKey))
	r.WithHeaders(h)
	return rawReq, nil
}
