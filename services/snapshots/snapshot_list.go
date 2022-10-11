package snapshots

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	snapListEndpoint = "/v1/projects/%s/snapshots"
)

type SnapshotListRequest struct {
	clo.Request
	ProjectID string
}

func (r *SnapshotListRequest) Make(ctx context.Context, cli *clo.ApiClient) (SnapshotListResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return SnapshotListResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return SnapshotListResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp SnapshotListResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return SnapshotListResponse{}, e
	}
	return resp, nil
}

func (r *SnapshotListRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(snapListEndpoint, r.ProjectID)
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

func (r *SnapshotListRequest) OrderBy(of string) *SnapshotListRequest {
	r.WithQueryParams(map[string][]string{"order": {of}})
	return r
}

type FilteringField struct {
	FieldName string
	Condition string
	Value     string
}

func (r *SnapshotListRequest) FilterBy(ff FilteringField) *SnapshotListRequest {
	switch ff.Condition {
	case "gt", "gte", "lt", "lte", "range", "in":
		condString := fmt.Sprintf("%s__%s", ff.FieldName, ff.Condition)
		r.WithQueryParams(map[string][]string{condString: {ff.Value}})
	}
	return r
}

type SnapshotPaginatorOptions struct {
	Limit  int
	Offset int
}

type SnapshotListPaginator struct {
	op       SnapshotPaginatorOptions
	client   *clo.ApiClient
	params   SnapshotListRequest
	lastPage bool
}

func (lp *SnapshotListPaginator) LastPage() bool {
	return lp.lastPage
}

func NewSnapshotListPaginator(client *clo.ApiClient, params SnapshotListRequest, op SnapshotPaginatorOptions) (*SnapshotListPaginator, error) {
	if op.Limit == 0 {
		return nil, fmt.Errorf("op.Limit should not be 0")
	}
	lp := SnapshotListPaginator{
		client: client,
		params: params,
		op:     op,
	}
	return &lp, nil
}

func (lp *SnapshotListPaginator) NextPage(ctx context.Context) (SnapshotListResponse, error) {
	if lp.LastPage() {
		return SnapshotListResponse{}, fmt.Errorf("no more pages")
	}
	lp.params.WithQueryParams(map[string][]string{"limit": {fmt.Sprintf("%d", lp.op.Limit)}})
	lp.params.WithQueryParams(map[string][]string{"offset": {fmt.Sprintf("%d", lp.op.Offset)}})
	r, e := lp.params.Make(ctx, lp.client)
	if e != nil {
		return SnapshotListResponse{}, e
	}
	lp.op.Offset += lp.op.Limit
	if r.Count <= lp.op.Limit+lp.op.Offset {
		lp.lastPage = true
	}
	return r, e
}
