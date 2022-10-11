package recipe

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	recipeListEndpoint = "/v1/recipes"
)

type RecipeListRequest struct {
	clo.Request
}

func (r *RecipeListRequest) Make(ctx context.Context, cli *clo.ApiClient) (RecipeListResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return RecipeListResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return RecipeListResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp RecipeListResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return RecipeListResponse{}, e
	}
	return resp, nil
}

func (r *RecipeListRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += recipeListEndpoint
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

func (r *RecipeListRequest) OrderBy(of string) *RecipeListRequest {
	r.WithQueryParams(map[string][]string{"order": {of}})
	return r
}

func (r *RecipeListRequest) FilterBy(ff FilteringField) *RecipeListRequest {
	switch ff.Condition {
	case "gt", "gte", "lt", "lte", "range", "in":
		condString := fmt.Sprintf("%s__%s", ff.FieldName, ff.Condition)
		r.WithQueryParams(map[string][]string{condString: {ff.Value}})
	}
	return r
}

type FilteringField struct {
	FieldName string
	Condition string
	Value     string
}

type RecipeListPaginator struct {
	op       PaginatorOptions
	client   *clo.ApiClient
	params   RecipeListRequest
	lastPage bool
}

type PaginatorOptions struct {
	Limit  int
	Offset int
}

func (lp *RecipeListPaginator) LastPage() bool {
	return lp.lastPage
}

func NewRecipeListPaginator(client *clo.ApiClient, params RecipeListRequest, op PaginatorOptions) (*RecipeListPaginator, error) {
	if op.Limit == 0 {
		return nil, fmt.Errorf("op.Limit should not be 0")
	}
	lp := RecipeListPaginator{
		client: client,
		params: params,
		op:     op,
	}
	return &lp, nil
}

func (lp *RecipeListPaginator) NextPage(ctx context.Context) (RecipeListResponse, error) {
	if lp.LastPage() {
		return RecipeListResponse{}, fmt.Errorf("no more pages")
	}
	lp.params.WithQueryParams(map[string][]string{"limit": {fmt.Sprintf("%d", lp.op.Limit)}})
	lp.params.WithQueryParams(map[string][]string{"offset": {fmt.Sprintf("%d", lp.op.Offset)}})
	r, e := lp.params.Make(ctx, lp.client)
	if e != nil {
		return RecipeListResponse{}, e
	}
	lp.op.Offset += lp.op.Limit
	if r.Count <= lp.op.Limit+lp.op.Offset {
		lp.lastPage = true
	}
	return r, e
}
