package clo

import (
	"encoding/json"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo/request_tools"
	"io"
	"net/http"
	"time"
)

type ApiClient struct {
	HttpClient HttpClient
	Log        Logger
	Options    map[string]interface{}
}

func NewDefaultClient(authKey string, baseUrl string) (*ApiClient, error) {
	if len(authKey) == 0 {
		return nil, fmt.Errorf("authKey should not be empty")
	}
	if len(baseUrl) == 0 {
		return nil, fmt.Errorf("baseUrl should not be empty")
	}
	s := ApiClient{HttpClient: http.DefaultClient}
	s.Options = map[string]interface{}{
		"auth_key": authKey,
		"base_url": baseUrl,
	}
	return &s, nil
}

func NewDefaultClientFromConfig(cfg Config) (*ApiClient, error) {
	if len(cfg.BaseUrl) == 0 {
		return nil, fmt.Errorf("Config.BaseUrl should be provided")
	}
	if len(cfg.AuthKey) == 0 {
		return nil, fmt.Errorf("Config.AuthKey should be provided")
	}
	cli := http.DefaultClient
	cli.Timeout = time.Duration(cfg.HttpTimeoutSeconds) * time.Second
	s := ApiClient{
		HttpClient: http.DefaultClient,
		Options:    cfg.ToMap(),
	}
	return &s, nil
}

func (cli *ApiClient) MakeRequest(req *http.Request) (*http.Response, error) {
	resp, e := cli.HttpClient.Do(req)
	if cli.Log != nil {
		cli.Log.Tracef("resp: %v, err: %v\n", resp, e)
	}
	if e != nil {
		return nil, e
	}
	if request_tools.IsError(resp.StatusCode) {
		var de request_tools.DefaultError
		defer resp.Body.Close()
		e = json.NewDecoder(resp.Body).Decode(&de)
		switch {
		case e == io.EOF:
			return nil, fmt.Errorf("error with an empty body, status code is : %d", resp.StatusCode)
		case e != nil:
			return nil, fmt.Errorf("can't decode an error body: %v", resp.Body)
		}
		return nil, de
	}
	return resp, nil
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
