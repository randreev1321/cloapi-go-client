package disks

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	localListEndpoint = "/v1/projects/%s/local-disks"
)

type LocalListRequest struct {
	clo.Request
	ProjectID string
}

func (r *LocalListRequest) Make(ctx context.Context, cli *clo.ApiClient) (LocalListResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return LocalListResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return LocalListResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp LocalListResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return LocalListResponse{}, e
	}
	return resp, nil
}

func (r *LocalListRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(localListEndpoint, r.ProjectID)
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

func (r *LocalListRequest) OrderBy(of string) *LocalListRequest {
	r.WithQueryParams(map[string][]string{"order": {of}})
	return r
}

func (r *LocalListRequest) FilterBy(ff FilteringField) *LocalListRequest {
	switch ff.Condition {
	case "gt", "gte", "lt", "lte", "range", "in":
		condString := fmt.Sprintf("%s__%s", ff.FieldName, ff.Condition)
		r.WithQueryParams(map[string][]string{condString: {ff.Value}})
	}
	return r
}

type LocalListPaginator struct {
	op       PaginatorOptions
	client   *clo.ApiClient
	params   *LocalListRequest
	lastPage bool
}

func (lp *LocalListPaginator) LastPage() bool {
	return lp.lastPage
}

func NewLocalListPaginator(client *clo.ApiClient, params *LocalListRequest, op PaginatorOptions) (*LocalListPaginator, error) {
	if op.Limit == 0 {
		return nil, fmt.Errorf("op.Limit should not be 0")
	}
	lp := LocalListPaginator{
		client: client,
		params: params,
		op:     op,
	}
	return &lp, nil
}

func (lp *LocalListPaginator) NextPage(ctx context.Context) (LocalListResponse, error) {
	if lp.LastPage() {
		return LocalListResponse{}, fmt.Errorf("no more pages")
	}
	lp.params.WithQueryParams(map[string][]string{"limit": {fmt.Sprintf("%d", lp.op.Limit)}})
	lp.params.WithQueryParams(map[string][]string{"offset": {fmt.Sprintf("%d", lp.op.Offset)}})
	r, e := lp.params.Make(ctx, lp.client)
	if e != nil {
		return LocalListResponse{}, e
	}
	lp.op.Offset += lp.op.Limit
	if r.Count <= lp.op.Limit+lp.op.Offset {
		lp.lastPage = true
	}
	return r, e
}
