package snapshots

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	snapDetailEndpoint = "/v1/snapshots/%s/detail"
)

type SnapshotDetailRequest struct {
	clo.Request
	SnapshotID string
}

func (r *SnapshotDetailRequest) Make(ctx context.Context, cli *clo.ApiClient) (SnapshotDetailResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return SnapshotDetailResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return SnapshotDetailResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp SnapshotDetailResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return SnapshotDetailResponse{}, e
	}
	return resp, nil
}

func (r *SnapshotDetailRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(snapDetailEndpoint, r.SnapshotID)
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
