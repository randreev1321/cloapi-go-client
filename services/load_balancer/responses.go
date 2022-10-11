package load_balancer

import (
	"encoding/json"
	"github.com/clo-ru/cloapi-go-client/clo"
	"io"
)

type BalancerListResponse struct {
	clo.Response
	Count   int                  `json:"count"`
	Results []BalancerDetailItem `json:"results"`
}

func (r *BalancerListResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type BalancerDetailResponse struct {
	clo.Response
	Result BalancerDetailItem `json:"result"`
}

func (r *BalancerDetailResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type BalancerCreateResponse struct {
	clo.Response
	Result BalancerDetailItem `json:"result"`
}

func (r *BalancerCreateResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type BalancerDetailItem struct {
	ID                string                 `json:"id"`
	ProjectID         string                 `json:"project_id"`
	Name              string                 `json:"name"`
	Algorithm         string                 `json:"algorithm"`
	Status            string                 `json:"status"`
	CreatedIn         string                 `json:"created_in"`
	UpdatedIn         string                 `json:"updated_in"`
	RulesCount        int                    `json:"rules_count"`
	RulesLimit        int                    `json:"rules_limit"`
	SessionPersistent bool                   `json:"session_persistent"`
	HealthMonitor     BalancerMonitorDetails `json:"healthmonitor"`
	Addresses         []BalancerAddress      `json:"addresses"`
	Links             []clo.Link             `json:"links"`
}

type BalancerMonitorDetails struct {
	Type          string `json:"type"`
	UrlPath       string `json:"url_path"`
	HttpMethod    string `json:"http_method"`
	ExpectedCodes string `json:"expected_codes"`
	Delay         int    `json:"delay"`
	Timeout       int    `json:"timeout"`
	MaxRetries    int    `json:"max_retries"`
}

type BalancerAddress struct {
	ID             string `json:"id"`
	Ptr            string `json:"ptr"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	MacAddr        string `json:"mac_addr"`
	Version        int    `json:"version"`
	External       bool   `json:"external"`
	DdosProtection bool   `json:"ddos_protection"`
}

type BalancerRulesListResponse struct {
	clo.Response
	Count   int                       `json:"count"`
	Results []BalancerRulesDetailItem `json:"results"`
}

func (r *BalancerRulesListResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type BalancerRulesDetailItem struct {
	ID                   string                  `json:"id"`
	ExternalProtocolPort int                     `json:"external_protocol_port"`
	InternalProtocolPort int                     `json:"internal_protocol_port"`
	Server               AttachedToEntityDetails `json:"server"`
	Loadbalancer         AttachedToEntityDetails `json:"loadbalancer"`
	Links                []clo.Link              `json:"links"`
}

type AttachedToEntityDetails struct {
	ID    string     `json:"id"`
	Links []clo.Link `json:"links"`
}
