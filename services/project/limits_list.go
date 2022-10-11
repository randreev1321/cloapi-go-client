package project

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	limitsListEndpoint = "/v1/projects/%s/limits"
)

type LimitsListRequest struct {
	clo.Request
	ProjectID string
}

func (r *LimitsListRequest) Make(ctx context.Context, cli *clo.ApiClient) (LimitsListResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return LimitsListResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return LimitsListResponse{}, requestError
	}
	var resp LimitsListResponse
	defer rawResp.Body.Close()
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return LimitsListResponse{}, e
	}
	return resp, nil
}

func (r *LimitsListRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(limitsListEndpoint, r.ProjectID)
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

func (r *LimitsListRequest) OrderBy(of string) *LimitsListRequest {
	r.WithQueryParams(map[string][]string{"order": {of}})
	return r
}

func (r *LimitsListRequest) FilterBy(ff FilteringField) *LimitsListRequest {
	switch ff.Condition {
	case "gt", "gte", "lt", "lte", "range", "in":
		condString := fmt.Sprintf("%s__%s", ff.FieldName, ff.Condition)
		r.WithQueryParams(map[string][]string{condString: {ff.Value}})
	}
	return r
}

type LimitsListPaginator struct {
	op       PaginatorOptions
	client   *clo.ApiClient
	params   LimitsListRequest
	lastPage bool
}

func (lp *LimitsListPaginator) LastPage() bool {
	return lp.lastPage
}

func NewLimitsListPaginator(client *clo.ApiClient, params LimitsListRequest, op PaginatorOptions) (*LimitsListPaginator, error) {
	if op.Limit == 0 {
		return nil, fmt.Errorf("op.Limit should not be 0")
	}
	lp := LimitsListPaginator{
		client: client,
		params: params,
		op:     op,
	}
	return &lp, nil
}

func (lp *LimitsListPaginator) NextPage(ctx context.Context) (LimitsListResponse, error) {
	if lp.LastPage() {
		return LimitsListResponse{}, fmt.Errorf("no more pages")
	}
	lp.params.WithQueryParams(map[string][]string{"limit": {fmt.Sprintf("%d", lp.op.Limit)}})
	lp.params.WithQueryParams(map[string][]string{"offset": {fmt.Sprintf("%d", lp.op.Offset)}})
	r, e := lp.params.Make(ctx, lp.client)
	if e != nil {
		return LimitsListResponse{}, e
	}
	lp.op.Offset += lp.op.Limit
	if r.Count <= lp.op.Limit+lp.op.Offset {
		lp.lastPage = true
	}
	return r, e
}
