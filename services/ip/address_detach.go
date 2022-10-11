package ip

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	addressDetachEndpoint = "/v1/addresses/%s/detach"
)

type AddressDetachRequest struct {
	clo.Request
	AddressID string
}

func (r *AddressDetachRequest) Make(ctx context.Context, cli *clo.ApiClient) (AddressDetachResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return AddressDetachResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return AddressDetachResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp AddressDetachResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return AddressDetachResponse{}, e
	}
	return resp, nil
}

func (r *AddressDetachRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(addressDetachEndpoint, r.AddressID)
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
