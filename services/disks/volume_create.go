package disks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	volumeCreateEndpoint = "/v1/projects/%s/volumes"
)

type VolumeCreateRequest struct {
	clo.Request
	ProjectID string
	Body      VolumeCreateBody
}

type VolumeCreateBody struct {
	Name       string `json:"name"`
	Size       int    `json:"size"`
	Autorename bool   `json:"autorename"`
}

func (r *VolumeCreateRequest) Make(ctx context.Context, cli *clo.ApiClient) (VolumeCreateResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return VolumeCreateResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return VolumeCreateResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp VolumeCreateResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return VolumeCreateResponse{}, e
	}
	return resp, nil
}

func (r *VolumeCreateRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(volumeCreateEndpoint, r.ProjectID)
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
