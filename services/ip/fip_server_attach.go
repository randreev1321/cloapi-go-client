package ip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	fipServerAttachEndpoint = "/v1/floatingips/%s/attach/server"
)

type FipServerAttachRequest struct {
	clo.Request
	FloatingID string
	Body       FipAttachBody
}

type FipAttachBody struct {
	ServerID string `json:"id"`
}

func (r *FipServerAttachRequest) Make(ctx context.Context, cli *clo.ApiClient) (FipServerAttachResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return FipServerAttachResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return FipServerAttachResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp FipServerAttachResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return FipServerAttachResponse{}, e
	}
	return resp, nil
}

func (r *FipServerAttachRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(fipServerAttachEndpoint, r.FloatingID)
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
