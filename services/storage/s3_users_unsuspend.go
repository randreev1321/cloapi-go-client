package storage

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	s3UserUnsuspendEndpoint = "/v1/s3_users/%s/unsuspend"
)

type S3UserUnsuspendRequest struct {
	clo.Request
	UserID string
}

func (r *S3UserUnsuspendRequest) Make(ctx context.Context, cli *clo.ApiClient) (S3UserUnsuspendResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return S3UserUnsuspendResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return S3UserUnsuspendResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp S3UserUnsuspendResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return S3UserUnsuspendResponse{}, e
	}
	return resp, nil
}

func (r *S3UserUnsuspendRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(s3UserUnsuspendEndpoint, r.UserID)
	rawReq, e := http.NewRequestWithContext(
		ctx, http.MethodPost, baseUrl, nil,
	)
	if e != nil {
		return nil, e
	}
	h := http.Header{}
	h.Add("Authorization", fmt.Sprintf("Bearer %s", authKey))
	r.WithHeaders(h)
	return rawReq, nil
}
