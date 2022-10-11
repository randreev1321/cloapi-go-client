package project

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	imageListEndpoint = "/v1/projects/%s/images"
)

type ImageListRequest struct {
	clo.Request
	ProjectID string
}

func (r *ImageListRequest) Make(ctx context.Context, cli *clo.ApiClient) (ImageListResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return ImageListResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return ImageListResponse{}, requestError
	}
	var resp ImageListResponse
	defer rawResp.Body.Close()
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return ImageListResponse{}, e
	}
	return resp, nil
}

func (r *ImageListRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(imageListEndpoint, r.ProjectID)
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

func (r *ImageListRequest) OrderBy(of string) *ImageListRequest {
	r.WithQueryParams(map[string][]string{"order": {of}})
	return r
}

func (r *ImageListRequest) FilterBy(ff FilteringField) *ImageListRequest {
	switch ff.Condition {
	case "gt", "gte", "lt", "lte", "range", "in":
		condString := fmt.Sprintf("%s__%s", ff.FieldName, ff.Condition)
		r.WithQueryParams(map[string][]string{condString: {ff.Value}})
	}
	return r
}

type ImageListPaginator struct {
	op       PaginatorOptions
	client   *clo.ApiClient
	params   ImageListRequest
	lastPage bool
}

func (lp *ImageListPaginator) LastPage() bool {
	return lp.lastPage
}

func NewImageListPaginator(client *clo.ApiClient, params ImageListRequest, op PaginatorOptions) (*ImageListPaginator, error) {
	if op.Limit == 0 {
		return nil, fmt.Errorf("op.Limit should not be 0")
	}
	lp := ImageListPaginator{
		client: client,
		params: params,
		op:     op,
	}
	return &lp, nil
}

func (lp *ImageListPaginator) NextPage(ctx context.Context) (ImageListResponse, error) {
	if lp.LastPage() {
		return ImageListResponse{}, fmt.Errorf("no more pages")
	}
	lp.params.WithQueryParams(map[string][]string{"limit": {fmt.Sprintf("%d", lp.op.Limit)}})
	lp.params.WithQueryParams(map[string][]string{"offset": {fmt.Sprintf("%d", lp.op.Offset)}})
	r, e := lp.params.Make(ctx, lp.client)
	if e != nil {
		return ImageListResponse{}, e
	}
	lp.op.Offset += lp.op.Limit
	if r.Count <= lp.op.Limit+lp.op.Offset {
		lp.lastPage = true
	}
	return r, e
}
