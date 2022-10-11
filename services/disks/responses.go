package disks

import (
	"encoding/json"
	"github.com/clo-ru/cloapi-go-client/clo"
	"io"
)

type LocalListResponse struct {
	clo.Response
	Count   int            `json:"count"`
	Results []ResponseItem `json:"results"`
}

func (r *LocalListResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type LocalDetailResponse struct {
	clo.Response
	Result ResponseItem `json:"result"`
}

func (r *LocalDetailResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type VolumeListResponse struct {
	clo.Response
	Count   int            `json:"count"`
	Results []VolumeDetail `json:"results"`
}

func (r *VolumeListResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type VolumeDetailResponse struct {
	clo.Response
	Result VolumeDetail `json:"result"`
}

func (r *VolumeDetailResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type VolumeCreateResponse struct {
	clo.Response
	Result VolumeDetail `json:"result"`
}

func (r *VolumeCreateResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type VolumeAttachResponse struct {
	clo.Response
	Result VolumeAttachItem `json:"result"`
}

func (r *VolumeAttachResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type VolumeAttachItem struct {
	Device     string             `json:"device"`
	MountCmd   string             `json:"mount_cmd"`
	Mountpoint string             `json:"mountpoint"`
	Volume     VolumeAttachDetail `json:"volume"`
	Server     AttachedToServer   `json:"server"`
}

type VolumeAttachDetail struct {
	ID    string     `json:"id"`
	Links []clo.Link `json:"links"`
}

type VolumeDetail struct {
	ResponseItem
	Undetachable bool   `json:"undetachable"`
	Device       string `json:"device"`
	Description  string `json:"description"`
}

type ResponseItem struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Status           string           `json:"status"`
	CreatedIn        string           `json:"created_in"`
	Size             int              `json:"size"`
	Bootable         bool             `json:"bootable"`
	AttachedToServer AttachedToServer `json:"attached_to_server"`
	Links            []clo.Link       `json:"links"`
}

type AttachedToServer struct {
	ID    string     `json:"id"`
	Links []clo.Link `json:"links"`
}
