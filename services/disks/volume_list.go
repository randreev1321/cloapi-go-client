package disks

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	volumeListEndpoint = "/v1/projects/%s/volumes"
)

type VolumeListRequest struct {
	clo.Request
	ProjectID string
}

func (r *VolumeListRequest) Make(ctx context.Context, cli *clo.ApiClient) (VolumeListResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return VolumeListResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return VolumeListResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp VolumeListResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return VolumeListResponse{}, e
	}
	return resp, nil
}

func (r *VolumeListRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(volumeListEndpoint, r.ProjectID)
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

func (r *VolumeListRequest) OrderBy(of string) *VolumeListRequest {
	r.WithQueryParams(map[string][]string{"order": {of}})
	return r
}

func (r *VolumeListRequest) FilterBy(ff FilteringField) *VolumeListRequest {
	switch ff.Condition {
	case "gt", "gte", "lt", "lte", "range", "in":
		condString := fmt.Sprintf("%s__%s", ff.FieldName, ff.Condition)
		r.WithQueryParams(map[string][]string{condString: {ff.Value}})
	}
	return r
}

type VolumeListPaginator struct {
	op       PaginatorOptions
	client   *clo.ApiClient
	params   VolumeListRequest
	lastPage bool
}

func (lp *VolumeListPaginator) LastPage() bool {
	return lp.lastPage
}

func NewVolumeListPaginator(client *clo.ApiClient, params VolumeListRequest, op PaginatorOptions) (*VolumeListPaginator, error) {
	if op.Limit == 0 {
		return nil, fmt.Errorf("op.Limit should not be 0")
	}
	lp := VolumeListPaginator{
		client: client,
		params: params,
		op:     op,
	}
	return &lp, nil
}

func (lp *VolumeListPaginator) NextPage(ctx context.Context) (VolumeListResponse, error) {
	if lp.LastPage() {
		return VolumeListResponse{}, fmt.Errorf("no more pages")
	}
	lp.params.WithQueryParams(map[string][]string{"limit": {fmt.Sprintf("%d", lp.op.Limit)}})
	lp.params.WithQueryParams(map[string][]string{"offset": {fmt.Sprintf("%d", lp.op.Offset)}})
	r, e := lp.params.Make(ctx, lp.client)
	if e != nil {
		return VolumeListResponse{}, e
	}
	lp.op.Offset += lp.op.Limit
	if r.Count <= lp.op.Limit+lp.op.Offset {
		lp.lastPage = true
	}
	return r, e
}
