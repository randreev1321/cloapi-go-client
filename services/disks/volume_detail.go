package disks

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"net/http"
)

const (
	volumeDetailEndpoint = "/v1/volumes/%s/detail"
)

type VolumeDetailRequest struct {
	clo.Request
	VolumeID string
}

func (r *VolumeDetailRequest) Make(ctx context.Context, cli *clo.ApiClient) (VolumeDetailResponse, error) {
	rawReq, e := r.buildRequest(ctx, cli.Options)
	if e != nil {
		return VolumeDetailResponse{}, e
	}
	rawResp, requestError := r.MakeRequest(rawReq, cli)
	if requestError != nil {
		return VolumeDetailResponse{}, requestError
	}
	defer rawResp.Body.Close()
	var resp VolumeDetailResponse
	if e = resp.FromJsonBody(rawResp.Body); e != nil {
		return VolumeDetailResponse{}, e
	}
	return resp, nil
}

func (r *VolumeDetailRequest) buildRequest(ctx context.Context, cliOptions map[string]interface{}) (*http.Request, error) {
	authKey, ok := cliOptions["auth_key"].(string)
	if !ok {
		return nil, fmt.Errorf("auth_key client options should be a string, %T got", authKey)
	}
	baseUrl, ok := cliOptions["base_url"].(string)
	if !ok {
		return nil, fmt.Errorf("base_url client options should be a string, %T got", baseUrl)
	}
	baseUrl += fmt.Sprintf(volumeDetailEndpoint, r.VolumeID)
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
