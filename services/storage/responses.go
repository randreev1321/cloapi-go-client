package storage

import (
	"encoding/json"
	"github.com/clo-ru/cloapi-go-client/clo"
	"io"
)

type S3KeysResetResponse struct {
	clo.Response
	Result S3KeysResponse `json:"result"`
}

func (r *S3KeysResetResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type S3KeysGetResponse struct {
	clo.Response
	Result S3KeysResponse `json:"result"`
}

func (r *S3KeysGetResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type S3KeysResponse struct {
	User      string `json:"user"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type S3UserListResponse struct {
	clo.Response
	Count   int            `json:"count"`
	Results []ResponseItem `json:"results"`
}

func (r *S3UserListResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type S3UserPatchResponse ResponseItemResult

func (r *S3UserPatchResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type S3UserCreateResponse ResponseItemResult

func (r *S3UserCreateResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type S3UserDetailResponse ResponseItemResult

func (r *S3UserDetailResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type S3UserSuspendResponse ResponseItemResult

func (r *S3UserSuspendResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type S3UserUnsuspendResponse ResponseItemResult

func (r *S3UserUnsuspendResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type S3UserQuotaPatchResponse ResponseItemResult

func (r *S3UserQuotaPatchResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type ResponseItemResult struct {
	clo.Response
	Result ResponseItem `json:"result"`
}

type ResponseItem struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	CanonicalName string      `json:"canonical_name"`
	Status        string      `json:"status"`
	Tenant        string      `json:"tenant"`
	MaxBuckets    int         `json:"max_buckets"`
	Quotas        []QuotaInfo `json:"quotas"`
}

type QuotaInfo struct {
	MaxSize    int    `json:"max_size"`
	MaxObjects int    `json:"max_objects"`
	Type       string `json:"type"`
}
