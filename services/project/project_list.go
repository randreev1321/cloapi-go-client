package project

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	projectListEndpoint = "/v1/projects"
)

type ProjectListRequest struct {
	clo.Request
}

func (r *ProjectListRequest) Make(ctx context.Context, cli *clo.ApiClient) (ProjectListResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return ProjectListResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return ProjectListResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp ProjectListResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return ProjectListResponse{}, e
	}
	return resp, nil
}

func (r *ProjectListRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += projectListEndpoint
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

func (r *ProjectListRequest) OrderBy(of string) *ProjectListRequest {
	r.WithQueryParams(map[string][]string{"order": {of}})
	return r
}

func (r *ProjectListRequest) FilterBy(ff FilteringField) *ProjectListRequest {
	switch ff.Condition {
	case "gt", "gte", "lt", "lte", "range", "in":
		condString := fmt.Sprintf("%s__%s", ff.FieldName, ff.Condition)
		r.WithQueryParams(map[string][]string{condString: {ff.Value}})
	}
	return r
}

type ProjectListPaginator struct {
	op       PaginatorOptions
	client   *clo.ApiClient
	params   ProjectListRequest
	lastPage bool
}

func (lp *ProjectListPaginator) LastPage() bool {
	return lp.lastPage
}

func NewProjectListPaginator(client *clo.ApiClient, params ProjectListRequest, op PaginatorOptions) (*ProjectListPaginator, error) {
	if op.Limit == 0 {
		return nil, fmt.Errorf("op.Limit should not be 0")
	}
	lp := ProjectListPaginator{
		client: client,
		params: params,
		op:     op,
	}
	return &lp, nil
}

func (lp *ProjectListPaginator) NextPage(ctx context.Context) (ProjectListResponse, error) {
	if lp.LastPage() {
		return ProjectListResponse{}, fmt.Errorf("no more pages")
	}
	lp.params.WithQueryParams(map[string][]string{"limit": {fmt.Sprintf("%d", lp.op.Limit)}})
	lp.params.WithQueryParams(map[string][]string{"offset": {fmt.Sprintf("%d", lp.op.Offset)}})
	r, e := lp.params.Make(ctx, lp.client)
	if e != nil {
		return ProjectListResponse{}, e
	}
	lp.op.Offset += lp.op.Limit
	if r.Count <= lp.op.Limit+lp.op.Offset {
		lp.lastPage = true
	}
	return r, e
}
