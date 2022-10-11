package ip

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	addressPrimaryChangeEndpoint = "/v1/addresses/%s/primary"
)

type AddressPrimaryChangeRequest struct {
	clo.Request
	AddressID string
}

func (r *AddressPrimaryChangeRequest) Make(ctx context.Context, cli *clo.ApiClient) (AddressPrimaryChangeResponse, error) {
	var resp AddressPrimaryChangeResponse
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return resp, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return resp, requestError
	}
	defer rawResp.Body.Close()
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return resp, e
	}
	return resp, nil
}

func (r *AddressPrimaryChangeRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(addressPrimaryChangeEndpoint, r.AddressID)
	rawReq, e := http.NewRequestWithContext(
		ctx, http.MethodPost, baseUrl, nil,
	)
	h := http.Header{}
	h.Add("Authorization", fmt.Sprintf("Bearer %s", authKey))
	r.WithHeaders(h)
	if e != nil {
		return nil, e
	}
	return rawReq, nil
}
