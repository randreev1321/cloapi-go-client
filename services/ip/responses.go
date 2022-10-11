package ip

import (
	"encoding/json"
	"github.com/clo-ru/cloapi-go-client/clo"
	"io"
)

type FipListResponse struct {
	clo.Response
	Count   int         `json:"count"`
	Results []FipDetail `json:"results"`
}

func (r *FipListResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type FipCreateResponse struct {
	clo.Response
	Result CreateItem `json:"result"`
}

func (r *FipCreateResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type CreateItem struct {
	ID                string     `json:"id"`
	FloatingIPAddress string     `json:"floating_ip_address"`
	Links             []clo.Link `json:"links"`
}

type FipBalancerAttachResponse struct {
	clo.Response
	Result FipDetail `json:"result"`
}

func (r *FipBalancerAttachResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type FipServerAttachResponse struct {
	clo.Response
	Result FipDetail `json:"result"`
}

func (r *FipServerAttachResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type FipDetailResponse struct {
	clo.Response
	Result FipDetail `json:"result"`
}

func (r *FipDetailResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type FipDetail struct {
	ID                 string                  `json:"id"`
	Status             string                  `json:"status"`
	Ptr                string                  `json:"ptr"`
	CreatedIn          string                  `json:"created_in"`
	UpdatedIn          string                  `json:"updated_in"`
	FixedIpAddress     string                  `json:"fixed_ip_address"`
	FloatingIpAddress  string                  `json:"floating_ip_address"`
	DdosProtection     bool                    `json:"ddos_protection"`
	AttachedToServer   AttachedToEntityDetails `json:"attached_to_server"`
	AttachedToBalancer AttachedToEntityDetails `json:"attached_to_loadbalancer"`
	Links              []clo.Link              `json:"links"`
}
type AddressListResponse struct {
	clo.Response
	Count   int             `json:"count"`
	Results []AddressDetail `json:"results"`
}

func (r *AddressListResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type AddressCreateResponse AddressDetailResponse

func (r *AddressCreateResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type AddressPrimaryChangeResponse AddressDetailResponse

func (r *AddressPrimaryChangeResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type AddressDetachResponse AddressDetailResponse

func (r *AddressDetachResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type AddressAttachResponse AddressDetailResponse

func (r *AddressAttachResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type AddressDetailResponse struct {
	clo.Response
	Result AddressDetail `json:"result"`
}

func (r *AddressDetailResponse) FromJsonBody(body io.ReadCloser) error {
	if e := json.NewDecoder(body).Decode(r); e != nil {
		return e
	}
	return nil
}

type AddressDetail struct {
	ID             string            `json:"id"`
	Type           string            `json:"type"`
	Status         string            `json:"status"`
	Address        string            `json:"address"`
	Ptr            string            `json:"ptr"`
	CreatedIn      string            `json:"created_in"`
	UpdatedIn      string            `json:"updated_in"`
	IsPrimary      bool              `json:"is_primary"`
	DdosProtection bool              `json:"ddos_protection"`
	AttachedTo     AttachedToDetails `json:"attached_to"`
	Links          []clo.Link        `json:"links"`
}

type AttachedToDetails struct {
	ID     string     `json:"id"`
	Entity string     `json:"entity"`
	Links  []clo.Link `json:"links"`
}

type AttachedToEntityDetails struct {
	ID    string     `json:"id"`
	Links []clo.Link `json:"links"`
}
