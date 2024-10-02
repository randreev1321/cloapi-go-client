package load_balancer

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/v2/clo"
	"net/http"
)

const (
	balancerRulesDetailEndpoint = "%s/v2/loadbalancers/rules/%s"
)

type BalancerRulesDetailRequest struct {
	clo.Request
	ObjectId string
}

func (r *BalancerRulesDetailRequest) Do(ctx context.Context, cli *clo.ApiClient) (*BalancerRuleDetailResponse, error) {
	res := &BalancerRuleDetailResponse{}
	return res, cli.DoRequest(ctx, r, res)
}

func (r *BalancerRulesDetailRequest) Build(ctx context.Context, baseUrl string, authToken string) (*http.Request, error) {
	return r.BuildRaw(ctx, http.MethodGet, fmt.Sprintf(balancerRulesDetailEndpoint, baseUrl, r.ObjectId), authToken, nil)
}
