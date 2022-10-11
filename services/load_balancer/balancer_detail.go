package load_balancer

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	balancerDetailEndpoint = "/v1/loadbalancers/%s/detail"
)

type BalancerDetailRequest struct {
	clo.Request
	ServerID string
}

func (r *BalancerDetailRequest) Make(ctx context.Context, cli *clo.ApiClient) (BalancerDetailResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return BalancerDetailResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return BalancerDetailResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp BalancerDetailResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return BalancerDetailResponse{}, e
	}
	return resp, nil
}

func (r *BalancerDetailRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(balancerDetailEndpoint, r.ServerID)
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
