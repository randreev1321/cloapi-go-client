package clo

import "encoding/json"

type Config struct {
	AuthKey            string `json:"auth_key"`
	BaseUrl            string `json:"base_url"`
	HttpTimeoutSeconds int    `json:"http_timeout_seconds"`
}

func (cfg *Config) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"auth_key": cfg.AuthKey,
		"base_url": cfg.BaseUrl,
		"timeout":  cfg.HttpTimeoutSeconds,
	}
}

func (cfg *Config) FromMap(opt map[string]interface{}) error {
	b, e := json.Marshal(opt)
	if e != nil {
		return e
	}
	if e := json.Unmarshal(b, cfg); e != nil {
		return e
	}
	return nil
}
