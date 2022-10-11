package sshkeys

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	keypairListEndpoint = "/v1/keypairs"
)

type KeyPairListRequest struct {
	clo.Request
}

func (r *KeyPairListRequest) Make(ctx context.Context, cli *clo.ApiClient) (KeyPairListResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return KeyPairListResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return KeyPairListResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp KeyPairListResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return KeyPairListResponse{}, e
	}
	return resp, nil
}

func (r *KeyPairListRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += keypairListEndpoint
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

func (r *KeyPairListRequest) OrderBy(of string) *KeyPairListRequest {
	r.WithQueryParams(map[string][]string{"order": {of}})
	return r
}

type FilteringField struct {
	FieldName string
	Condition string
	Value     string
}

func (r *KeyPairListRequest) FilterBy(ff FilteringField) *KeyPairListRequest {
	switch ff.Condition {
	case "gt", "gte", "lt", "lte", "range", "in":
		condString := fmt.Sprintf("%s__%s", ff.FieldName, ff.Condition)
		r.WithQueryParams(map[string][]string{condString: {ff.Value}})
	}
	return r
}

type KeyPairPaginatorOptions struct {
	Limit  int
	Offset int
}

type KeyPairListPaginator struct {
	op       KeyPairPaginatorOptions
	client   *clo.ApiClient
	params   KeyPairListRequest
	lastPage bool
}

func (lp *KeyPairListPaginator) LastPage() bool {
	return lp.lastPage
}

func NewKeyPairListPaginator(client *clo.ApiClient, params KeyPairListRequest, op KeyPairPaginatorOptions) (*KeyPairListPaginator, error) {
	if op.Limit == 0 {
		return nil, fmt.Errorf("op.Limit should not be 0")
	}
	lp := KeyPairListPaginator{
		client: client,
		params: params,
		op:     op,
	}
	return &lp, nil
}

func (lp *KeyPairListPaginator) NextPage(ctx context.Context) (KeyPairListResponse, error) {
	if lp.LastPage() {
		return KeyPairListResponse{}, fmt.Errorf("no more pages")
	}
	lp.params.WithQueryParams(map[string][]string{"limit": {fmt.Sprintf("%d", lp.op.Limit)}})
	lp.params.WithQueryParams(map[string][]string{"offset": {fmt.Sprintf("%d", lp.op.Offset)}})
	r, e := lp.params.Make(ctx, lp.client)
	if e != nil {
		return KeyPairListResponse{}, e
	}
	lp.op.Offset += lp.op.Limit
	if r.Count <= lp.op.Limit+lp.op.Offset {
		lp.lastPage = true
	}
	return r, e
}
