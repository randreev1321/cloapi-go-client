package sshkeys

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	keypairDetailEndpoint = "/v1/keypairs/%s/detail"
)

type KeyPairDetailRequest struct {
	clo.Request
	KeypairID string
}

func (r *KeyPairDetailRequest) Make(ctx context.Context, cli *clo.ApiClient) (KeyPairDetailResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return KeyPairDetailResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return KeyPairDetailResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp KeyPairDetailResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return KeyPairDetailResponse{}, e
	}
	return resp, nil
}

func (r *KeyPairDetailRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(keypairDetailEndpoint, r.KeypairID)
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
