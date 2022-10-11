package sshkeys

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	keypairCreateEndpoint = "/v1/keypairs"
)

type KeyPairCreateRequest struct {
	clo.Request
	Body KeyPairCreateBody
}

type KeyPairCreateBody struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

func (r *KeyPairCreateRequest) Make(ctx context.Context, cli *clo.ApiClient) (KeyPairCreateResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return KeyPairCreateResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return KeyPairCreateResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp KeyPairCreateResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return KeyPairCreateResponse{}, e
	}
	return resp, nil
}

func (r *KeyPairCreateRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += keypairCreateEndpoint
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
