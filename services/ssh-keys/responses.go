package sshkeys

import (
	"encoding/json"
	"github.com/clo-ru/cloapi-go-client/clo"
	"io"
)

type KeyPairListResponse struct {
	clo.Response
	Count   int                   `json:"count"`
	Results []KeyPairResponseItem `json:"results"`
}

func (r *KeyPairListResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type KeyPairDetailResponse struct {
	clo.Response
	Result KeyPairResponseItem `json:"result"`
}

func (r *KeyPairDetailResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type KeyPairCreateResponse struct {
	clo.Response
	Result KeyPairResponseItem `json:"result"`
}

func (r *KeyPairCreateResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type KeyPairResponseItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}
