package load_balancer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clo-ru/cloapi-go-client/v2/clo"
)

const (
	balancerRulesDeleteEndpoint = "%s/v2/loadbalancers/rules/%s"
)

type BalancerRulesDeleteRequest struct {
	clo.Request
	ObjectId string
}

func (r *BalancerRulesDeleteRequest) Do(ctx context.Context, cli *clo.ApiClient) error {
	return cli.DoRequest(ctx, r, nil)
}

func (r *BalancerRulesDeleteRequest) Build(ctx context.Context, baseUrl string, authToken string) (*http.Request, error) {
	return r.BuildRaw(ctx, http.MethodDelete, fmt.Sprintf(balancerRulesDeleteEndpoint, baseUrl, r.ObjectId), authToken, nil)
}
