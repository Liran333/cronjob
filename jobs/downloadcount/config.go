/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package downloadcount for config

package downloadcount

type Config struct {
	Spec            string `json:"spec" required:"true"`
	OriginalDataUrl string `json:"original_data_url" required:"true"`
}