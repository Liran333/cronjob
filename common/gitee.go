/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package common
package common

import (
	"fmt"
	"net/url"
)

type Gitee struct {
	Token    string `json:"token" required:"true"`
	Endpoint string `json:"endpoint" required:"true"`
}

func (g *Gitee) CommitUrl(owner, repo, path string) string {
	commitUrl := fmt.Sprintf("%s/api/v5/repos/%s/%s/commits?path=%s&page=1&per_page=1&access_token=%s",
		g.Endpoint,
		owner,
		repo,
		url.QueryEscape(path),
		g.Token,
	)
	return commitUrl
}
