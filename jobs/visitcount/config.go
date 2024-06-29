/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package visitcount for config

package visitcount

type Config struct {
	Spec            string `json:"spec" required:"true"`
	OriginalDataUrl string `json:"original_data_url" required:"true"`
}
