package snapshots

import (
	"encoding/json"
	"github.com/clo-ru/cloapi-go-client/clo"
	"io"
)

type SnapshotListResponse struct {
	clo.Response
	Count   int                  `json:"count"`
	Results []SnapshotDetailItem `json:"results"`
}

func (r *SnapshotListResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type SnapshotDetailResponse struct {
	clo.Response
	Result SnapshotDetailItem `json:"result"`
}

func (r *SnapshotDetailResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type SnapshotDetailItem struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	CreatedIn    string         `json:"created_in"`
	DeletedIn    string         `json:"deleted_in"`
	Status       string         `json:"status"`
	Size         int            `json:"size"`
	ParentServer ResponseItem   `json:"parent_server"`
	ChildServers []ResponseItem `json:"child_servers"`
	Links        []clo.Link     `json:"links"`
}

type ResponseItem struct {
	ID    string     `json:"id"`
	Links []clo.Link `json:"links"`
}

type SnapshotCreateResponse struct {
	clo.Response
	Result ResponseItem `json:"result"`
}

func (r *SnapshotCreateResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type SnapshotRestoreResponse struct {
	clo.Response
	Result SnapshotRestoreItem `json:"result"`
}

func (r *SnapshotRestoreResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type SnapshotRestoreItem struct {
	Name   string       `json:"name"`
	Server ResponseItem `json:"server"`
}
