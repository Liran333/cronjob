/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package modeldeploy
package modeldeploy

import "fmt"

type Config struct {
	IconUpload     IconUpload     `json:"icon_upload"`
	DeployLocation DeployLocation `json:"deploy_location"`
}

type IconUpload struct {
	Bucket   string `json:"bucket" required:"true"`
	EndPoint string `json:"endpoint" required:"true"`
}

type DeployLocation struct {
	Owner string `json:"owner" required:"true"`
	Repo  string `json:"repo" required:"true"`
	Path  string `json:"path" required:"true"`
	Ref   string `json:"ref" required:"true"`
}

func (d DeployLocation) ToString() string {
	return fmt.Sprintf("the local file [%s/%s/%s/%s]", d.Owner, d.Repo, d.Ref, d.Path)
}

func (cfg *Config) SetDefault() {
	if cfg.DeployLocation.Ref == "" {
		cfg.DeployLocation.Ref = "master"
	}
}
