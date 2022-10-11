package servers

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	serverDetailEndpoint = "/v1/servers/%s/detail"
)

type ServerDetailRequest struct {
	clo.Request
	ServerID string
}

func (r *ServerDetailRequest) Make(ctx context.Context, cli *clo.ApiClient) (ServerDetailResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return ServerDetailResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return ServerDetailResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp ServerDetailResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return ServerDetailResponse{}, e
	}
	return resp, nil
}

func (r *ServerDetailRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(serverDetailEndpoint, r.ServerID)
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
