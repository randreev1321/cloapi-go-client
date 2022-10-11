package project

import (
	"encoding/json"
	"github.com/clo-ru/cloapi-go-client/clo"
	"io"
)

type ProjectListResponse struct {
	clo.Response
	Count   int               `json:"count"`
	Results []ProjectListItem `json:"results"`
}

func (r *ProjectListResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type ProjectListItem struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	Status         string             `json:"status"`
	CreatedIn      string             `json:"created_in"`
	StoppingReason string             `json:"stopping_reason"`
	HasAbuse       bool               `json:"has_abuse"`
	Summary        ProjectItemSummary `json:"summary"`
}

type ProjectItemSummary struct {
	FloatingIps int `json:"floating_ips"`
	Networks    int `json:"networks"`
	Servers     int `json:"servers"`
	Volumes     int `json:"volumes"`
}

type ImageListResponse struct {
	clo.Response
	Count   int             `json:"count"`
	Results []ImageListItem `json:"results"`
}

func (r *ImageListResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type ImageListItem struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	OperationSystem OperationSystemItem `json:"operation_system"`
}

type OperationSystemItem struct {
	Distribution string `json:"distribution"`
	OsFamily     string `json:"os_family"`
	Version      string `json:"version"`
}

type LimitsListResponse struct {
	clo.Response
	Count   int              `json:"count"`
	Results []LimitsListItem `json:"results"`
}

func (r *LimitsListResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type LimitsListItem struct {
	Max     int    `json:"max"`
	Used    int    `json:"used"`
	Name    string `json:"name"`
	Measure string `json:"measure"`
}
