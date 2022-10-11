package load_balancer

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	balancerRulesListEndpoint = "/v1/projects/%s/loadbalancers/rules"
)

type BalancerRulesListRequest struct {
	clo.Request
	ProjectID string
}

func (r *BalancerRulesListRequest) Make(ctx context.Context, cli *clo.ApiClient) (BalancerRulesListResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return BalancerRulesListResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return BalancerRulesListResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp BalancerRulesListResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return BalancerRulesListResponse{}, e
	}
	return resp, nil
}

func (r *BalancerRulesListRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(balancerRulesListEndpoint, r.ProjectID)
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

func (r *BalancerRulesListRequest) OrderBy(of string) *BalancerRulesListRequest {
	r.WithQueryParams(map[string][]string{"order": {of}})
	return r
}

func (r *BalancerRulesListRequest) FilterBy(ff FilteringField) *BalancerRulesListRequest {
	switch ff.Condition {
	case "gt", "gte", "lt", "lte", "range", "in":
		condString := fmt.Sprintf("%s__%s", ff.FieldName, ff.Condition)
		r.WithQueryParams(map[string][]string{condString: {ff.Value}})
	}
	return r
}

type BalancerRulesListPaginator struct {
	op       PaginatorOptions
	client   *clo.ApiClient
	params   BalancerRulesListRequest
	lastPage bool
}

func (lp *BalancerRulesListPaginator) LastPage() bool {
	return lp.lastPage
}

func NewBalancerRulesListPaginator(client *clo.ApiClient, params BalancerRulesListRequest, op PaginatorOptions) (*BalancerRulesListPaginator, error) {
	if op.Limit == 0 {
		return nil, fmt.Errorf("op.Limit should not be 0")
	}
	lp := BalancerRulesListPaginator{
		client: client,
		params: params,
		op:     op,
	}
	return &lp, nil
}

func (lp *BalancerRulesListPaginator) NextPage(ctx context.Context) (BalancerRulesListResponse, error) {
	if lp.LastPage() {
		return BalancerRulesListResponse{}, fmt.Errorf("no more pages")
	}
	lp.params.WithQueryParams(map[string][]string{"limit": {fmt.Sprintf("%d", lp.op.Limit)}})
	lp.params.WithQueryParams(map[string][]string{"offset": {fmt.Sprintf("%d", lp.op.Offset)}})
	r, e := lp.params.Make(ctx, lp.client)
	if e != nil {
		return BalancerRulesListResponse{}, e
	}
	lp.op.Offset += lp.op.Limit
	if r.Count <= lp.op.Limit+lp.op.Offset {
		lp.lastPage = true
	}
	return r, e
}
