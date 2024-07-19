/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package common
package common

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

var obsClient *obs.ObsClient

type Obs struct {
	AccessKey string `json:"access_key"    required:"true"`
	SecretKey string `json:"secret_key"    required:"true"`
	Endpoint  string `json:"endpoint"      required:"true"`
}

func InitObs(cfg *Obs) error {
	cli, err := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	if err != nil {
		return err
	}

	obsClient = cli

	return nil
}

func ObsClient() *obs.ObsClient {
	return obsClient
}
