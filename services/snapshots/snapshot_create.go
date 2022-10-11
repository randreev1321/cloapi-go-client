package snapshots

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	snapCreateEndpoint = "/v1/servers/%s/snapshot"
)

type SnapshotCreateRequest struct {
	clo.Request
	ServerID string
	Body     SnapshotCreateBody
}

type SnapshotCreateBody struct {
	Name string `json:"name"`
}

func (r *SnapshotCreateRequest) Make(ctx context.Context, cli *clo.ApiClient) (SnapshotCreateResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return SnapshotCreateResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return SnapshotCreateResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp SnapshotCreateResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return SnapshotCreateResponse{}, e
	}
	return resp, nil
}

func (r *SnapshotCreateRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(snapCreateEndpoint, r.ServerID)
	b := new(bytes.Buffer)
	if e := json.NewEncoder(b).Encode(r.Body); e != nil {
		return nil, fmt.Errorf("can't encode body parameters, %s", e.Error())
	}
	rawReq, e := http.NewRequestWithContext(
		ctx, http.MethodPost, baseUrl, b,
	)
	if e != nil {
		return nil, e
	}
	h := http.Header{}
	h.Add("Authorization", fmt.Sprintf("Bearer %s", authKey))
	r.WithHeaders(h)
	return rawReq, nil
}
