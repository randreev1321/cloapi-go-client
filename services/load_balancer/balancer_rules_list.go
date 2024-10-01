package load_balancer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clo-ru/cloapi-go-client/v2/clo"
)

const (
	projectRulesListEndpoint  = "%s/v2/projects/%s/loadbalancers/rules"
	balancerRulesListEndpoint = "%s/v2/loadbalancers/%s/rules"
)

type BalancerRulesListRequest struct {
	clo.Request
	ProjectID  string
	BalancerID string
}

func (r *BalancerRulesListRequest) Do(ctx context.Context, cli *clo.ApiClient) (*BalancerRuleListResponse, error) {
	resp := &BalancerRuleListResponse{}
	return resp, cli.DoRequest(ctx, r, resp)
}
func (r *BalancerRulesListRequest) Build(ctx context.Context, baseUrl string, authToken string) (*http.Request, error) {
	if r.BalancerID != "" {
		return r.BuildRaw(ctx, http.MethodGet, fmt.Sprintf(balancerRulesListEndpoint, baseUrl, r.BalancerID), authToken, nil)
	}
	return r.BuildRaw(ctx, http.MethodGet, fmt.Sprintf(projectRulesListEndpoint, baseUrl, r.ProjectID), authToken, nil)
}
func (r *BalancerRulesListRequest) OrderBy(of string) {
	r.WithQueryParams(clo.QueryParam{"order": {of}})
}
func (r *BalancerRulesListRequest) FilterBy(ff clo.FilteringField) { clo.AddFilterToRequest(r, ff) }
