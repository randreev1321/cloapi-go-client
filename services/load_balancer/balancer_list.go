package load_balancer

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	balancerListEndpoint = "/v1/projects/%s/loadbalancers"
)

type BalancerListRequest struct {
	clo.Request
	ProjectID string
}

func (r *BalancerListRequest) Make(ctx context.Context, cli *clo.ApiClient) (BalancerListResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return BalancerListResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return BalancerListResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp BalancerListResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return BalancerListResponse{}, e
	}
	return resp, nil
}

func (r *BalancerListRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(balancerListEndpoint, r.ProjectID)
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

func (r *BalancerListRequest) OrderBy(of string) *BalancerListRequest {
	r.WithQueryParams(map[string][]string{"order": {of}})
	return r
}

func (r *BalancerListRequest) FilterBy(ff FilteringField) *BalancerListRequest {
	switch ff.Condition {
	case "gt", "gte", "lt", "lte", "range", "in":
		condString := fmt.Sprintf("%s__%s", ff.FieldName, ff.Condition)
		r.WithQueryParams(map[string][]string{condString: {ff.Value}})
	}
	return r
}

type BalancerListPaginator struct {
	op       PaginatorOptions
	client   *clo.ApiClient
	params   BalancerListRequest
	lastPage bool
}

func (lp *BalancerListPaginator) LastPage() bool {
	return lp.lastPage
}

func NewBalancerListPaginator(client *clo.ApiClient, params BalancerListRequest, op PaginatorOptions) (*BalancerListPaginator, error) {
	if op.Limit == 0 {
		return nil, fmt.Errorf("op.Limit should not be 0")
	}
	lp := BalancerListPaginator{
		client: client,
		params: params,
		op:     op,
	}
	return &lp, nil
}

func (lp *BalancerListPaginator) NextPage(ctx context.Context) (BalancerListResponse, error) {
	if lp.LastPage() {
		return BalancerListResponse{}, fmt.Errorf("no more pages")
	}
	lp.params.WithQueryParams(map[string][]string{"limit": {fmt.Sprintf("%d", lp.op.Limit)}})
	lp.params.WithQueryParams(map[string][]string{"offset": {fmt.Sprintf("%d", lp.op.Offset)}})
	r, e := lp.params.Make(ctx, lp.client)
	if e != nil {
		return BalancerListResponse{}, e
	}
	lp.op.Offset += lp.op.Limit
	if r.Count <= lp.op.Limit+lp.op.Offset {
		lp.lastPage = true
	}
	return r, e
}
