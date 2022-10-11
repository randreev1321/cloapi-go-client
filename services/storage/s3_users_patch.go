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
	s3UserPatchEndpoint = "/v1/s3_users/%s"
)

type S3UserPatchRequest struct {
	clo.Request
	UserID string
	Body   S3UserPatchBody
}

type S3UserPatchBody struct {
	Name string `json:"name"`
}

func (r *S3UserPatchRequest) Make(ctx context.Context, cli *clo.ApiClient) (S3UserPatchResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return S3UserPatchResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return S3UserPatchResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp S3UserPatchResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return S3UserPatchResponse{}, e
	}
	return resp, nil
}

func (r *S3UserPatchRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(s3UserPatchEndpoint, r.UserID)
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
