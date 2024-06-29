/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package config

package config

import (
	"os"

	"github.com/openmerlin/merlin-sdk/httpclient"
	"sigs.k8s.io/yaml"

	"github.com/openmerlin/cronjob/jobs/downloadcount"
	"github.com/openmerlin/cronjob/jobs/visitcount"
)

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

type configItems interface {
	ConfigItems() []interface{}
}

type Config struct {
	Merlin        httpclient.Config    `json:"merlin"`
	DownloadCount downloadcount.Config `json:"download_count"`
	VisitCount    visitcount.Config    `json:"visit_count"`
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Merlin,
	}
}

func (cfg *Config) SetDefault() {
}

// Validate validates the configuration.
func (cfg *Config) Validate() error {
	return CheckConfig(cfg, "")
}

func loadFromYaml(path string, cfg interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, cfg)
}

func LoadConfig(path string, cfg *Config, remove bool) error {
	if remove {
		defer os.Remove(path)
	}

	if err := loadFromYaml(path, cfg); err != nil {
		return err
	}

	SetDefault(cfg)

	return Validate(cfg)
}

func SetDefault(cfg interface{}) {
	if f, ok := cfg.(configSetDefault); ok {
		f.SetDefault()
	}

	if f, ok := cfg.(configItems); ok {
		items := f.ConfigItems()

		for i := range items {
			SetDefault(items[i])
		}
	}
}

func Validate(cfg interface{}) error {
	if f, ok := cfg.(configValidate); ok {
		if err := f.Validate(); err != nil {
			return err
		}
	}

	if f, ok := cfg.(configItems); ok {
		items := f.ConfigItems()

		for i := range items {
			if err := Validate(items[i]); err != nil {
				return err
			}
		}
	}

	return nil
}
