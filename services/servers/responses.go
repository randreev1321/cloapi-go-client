package servers

import (
	"encoding/json"
	"github.com/clo-ru/cloapi-go-client/clo"
	"io"
)

type ServerListResponse struct {
	clo.Response
	Count   int              `json:"count"`
	Results []ServerListItem `json:"results"`
}

func (r *ServerListResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type ServerDetailResponse struct {
	clo.Response
	Result ServerDetailItem `json:"result"`
}

func (r *ServerDetailResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type ServerListItem struct {
	ID        string           `json:"id"`
	ProjectID string           `json:"project_id"`
	Name      string           `json:"name"`
	Status    string           `json:"status"`
	Image     string           `json:"image"`
	Recipe    string           `json:"recipe"`
	CreatedIn string           `json:"created_in"`
	InRescue  string           `json:"in_rescue"`
	Flavor    ServerFlavorData `json:"flavor"`
	DiskData  []ServerDiskData `json:"disk_data"`
	Links     []clo.Link       `json:"links"`
}

type ServerDetailItem struct {
	ID         string                `json:"id"`
	Name       string                `json:"name"`
	Image      string                `json:"image"`
	Recipe     string                `json:"recipe"`
	Status     string                `json:"status"`
	CreatedIn  string                `json:"created_in"`
	ProjectID  string                `json:"project_id"`
	RescueMode string                `json:"rescue_mode"`
	GuestAgent bool                  `json:"guest_agent"`
	Flavor     ServerFlavorData      `json:"flavor"`
	DiskData   []ServerDiskData      `json:"disk_data"`
	Addresses  []ServerDetailAddress `json:"addresses"`
}

type ServerDetailAddress struct {
	ID             string `json:"id"`
	Ptr            string `json:"ptr"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	MacAddr        string `json:"mac_addr"`
	Version        int    `json:"version"`
	DdosProtection bool   `json:"ddos_protection"`
	External       bool   `json:"external"`
}

type ServerCreateResponse struct {
	clo.Response
	Result ServerCreateItem `json:"result"`
}

func (r *ServerCreateResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type ServerCreateItem struct {
	ID        string           `json:"id"`
	Image     string           `json:"image"`
	Name      string           `json:"name"`
	Status    string           `json:"status"`
	CreatedIn string           `json:"created_in"`
	Flavor    ServerFlavorData `json:"flavor"`
	DiskData  []ServerDiskData `json:"disk_data"`
	Links     []clo.Link       `json:"links"`
}

type ServerFlavorData struct {
	Ram   int `json:"ram"`
	Vcpus int `json:"vcpus"`
}

type ServerDiskData struct {
	ID          string     `json:"id"`
	StorageType string     `json:"storage_type"`
	Links       []clo.Link `json:"links"`
}

type LicenseList struct {
	Count int `json:"count"`
}
