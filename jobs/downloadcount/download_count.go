/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package downloadcount for download count update

package downloadcount

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/openmerlin/merlin-sdk/statistic"
	"github.com/openmerlin/merlin-sdk/statistic/api"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

const (
	JobDownloadCount = "download-count"
)

func NewDownloadJob(cfg *Config) *download {
	return &download{cfg: cfg}
}

type download struct {
	cfg *Config
}

func (d *download) Type() string {
	return JobDownloadCount
}

func (d *download) Run() error {
	downloadData, err := d.fetchDownloadCounts(d.cfg.OriginalDataUrl)
	if err != nil {
		return fmt.Errorf("error fetching download counts: %w", err)
	}

	for _, repo := range downloadData.Data {
		if err = d.updateRepo(repo.RepoID, repo.Download); err != nil {
			logrus.Errorf("Failed to update download counts for repo ID %s: %s", repo.RepoID, err)
		}
	}

	return nil
}

func (d *download) updateRepo(id string, count int) error {
	_, err := api.UpdateRepo(&statistic.UpdateModel{DownloadCount: count}, id)
	if err != nil {
		return xerrors.Errorf("fail to use internal api, error: %w", err)
	}
	return nil
}

type DownloadData struct {
	Code int `json:"code"`
	Data []struct {
		Name     string `json:"name"`
		Download int    `json:"download"`
		RepoID   string `json:"repo_id"`
	} `json:"data"`
}

func (d *download) fetchDownloadCounts(url string) (*DownloadData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, xerrors.Errorf("fail to fetch data online, error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("fail to read data, error: %w", err)
	}

	var data DownloadData
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, xerrors.Errorf("fail to Unmarshal data, error: %w", err)
	}
	return &data, nil
}