package storage

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	s3KeysGetEndpoint = "/v1/s3_users/%s/keys"
)

type S3KeysGetRequest struct {
	clo.Request
	UserID string
}

func (r *S3KeysGetRequest) Make(ctx context.Context, cli *clo.ApiClient) (S3KeysGetResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return S3KeysGetResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return S3KeysGetResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp S3KeysGetResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return S3KeysGetResponse{}, e
	}
	return resp, nil
}

func (r *S3KeysGetRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(s3KeysGetEndpoint, r.UserID)
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
