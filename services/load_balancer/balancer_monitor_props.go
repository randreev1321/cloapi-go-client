package load_balancer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	balancerChangeMonitorEndpoint = "/v1/loadbalancers/%s/healthmonitor"
)

type BalancerChangeMonitorRequest struct {
	clo.Request
	BalancerID string
	Body       BalancerChangeMonitorBody
}

type BalancerChangeMonitorBody struct {
	Delay         int    `json:"delay"`
	Timeout       int    `json:"timeout"`
	MaxRetries    int    `json:"max_retries"`
	Type          string `json:"type"`
	UrlPath       string `json:"url_path"`
	HttpMethod    string `json:"http_method"`
	ExpectedCodes string `json:"expected_codes"`
}

func (r *BalancerChangeMonitorRequest) Make(ctx context.Context, cli *clo.ApiClient) error {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return e
	}
	_, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return requestError
	}
	return nil
}

func (r *BalancerChangeMonitorRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(balancerChangeMonitorEndpoint, r.BalancerID)
	b := new(bytes.Buffer)
	if e := json.NewEncoder(b).Encode(r.Body); e != nil {
		return nil, fmt.Errorf("can't encode body parameters, %s", e.Error())
	}
	rawReq, e := http.NewRequestWithContext(
		ctx, http.MethodPut, baseUrl, b,
	)
	if e != nil {
		return nil, e
	}
	h := http.Header{}
	h.Add("Authorization", fmt.Sprintf("Bearer %s", authKey))
	r.WithHeaders(h)
	return rawReq, nil
}
