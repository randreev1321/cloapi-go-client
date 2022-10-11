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
	balancerCreateEndpoint = "/v1/projects/%s/loadbalancers"
)

type BalancerCreateRequest struct {
	clo.Request
	ProjectID string
	Body      BalancerCreateBody
}

type BalancerCreateBody struct {
	Name               string              `json:"name"`
	Algorithm          string              `json:"algorithm"`
	SessionPersistence bool                `json:"session_persistence"`
	FloatingIP         BalancerBodyAddress `json:"floating_ip,omitempty"`
	HealthMonitor      BalancerBodyMonitor `json:"healthmonitor"`
	Rules              []BalancerBodyRules `json:"rules"`
}

type BalancerBodyAddress struct {
	ID             string `json:"id"`
	DdosProtection bool   `json:"ddos_protection"`
}

type BalancerBodyMonitor struct {
	Delay         int    `json:"delay"`
	Timeout       int    `json:"timeout"`
	MaxRetries    int    `json:"max_retries"`
	Type          string `json:"type"`
	UrlPath       string `json:"url_path,omitempty"`
	HttpMethod    string `json:"http_method,omitempty"`
	ExpectedCodes string `json:"expected_codes,omitempty"`
}

type BalancerBodyRules struct {
	PortID               string `json:"port_id"`
	ExternalProtocolPort int    `json:"external_protocol_port"`
	InternalProtocolPort int    `json:"internal_protocol_port"`
}

func (r *BalancerCreateRequest) Make(ctx context.Context, cli *clo.ApiClient) (BalancerCreateResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return BalancerCreateResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return BalancerCreateResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp BalancerCreateResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return BalancerCreateResponse{}, e
	}
	return resp, nil
}

func (r *BalancerCreateRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(balancerCreateEndpoint, r.ProjectID)
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
