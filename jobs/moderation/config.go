/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package moderation for config

package moderation

type Config struct {
	InitPageNum    int `json:"init_page_num" required:"true"`
	InitLastTime   int `json:"init_last_time" required:"true"`
	PicPageNum     int `json:"pic_page_num" required:"true"`
	PicLastTime    int `json:"pic_last_time" required:"true"`
	DocPageNum     int `json:"doc_page_num" required:"true"`
	DocLastTime    int `json:"doc_last_time" required:"true"`
	ReadmePageNum  int `json:"readme_page_num" required:"true"`
	ReadmeLastTime int `json:"readme_last_time" required:"true"`
	VideoPageNum   int `json:"video_page_num" required:"true"`
	VideoLastTime  int `json:"video_last_time" required:"true"`
	AudioPageNum   int `json:"audio_page_num" required:"true"`
	AudioLastTime  int `json:"audio_last_time" required:"true"`
}
