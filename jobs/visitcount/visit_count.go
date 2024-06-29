/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package visitcount for

package visitcount

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
	JobVisitCount = "visit-count"
)

func NewVisitJob(cfg *Config) *visit {
	return &visit{cfg: cfg}
}

type visit struct {
	cfg *Config
}

func (d *visit) Type() string {
	return JobVisitCount
}

func (d *visit) Run() error {
	downloadData, err := d.fetchVisitCounts(d.cfg.OriginalDataUrl)
	if err != nil {
		return fmt.Errorf("error fetching download counts: %w", err)
	}

	for _, repo := range downloadData.Data {
		if err = d.updateRepo(repo.RepoID, repo.Visit); err != nil {
			logrus.Errorf("Failed to update download counts for repo ID %s: %s", repo.RepoID, err)
		}
	}

	return nil
}

func (d *visit) updateRepo(id string, count int) error {
	_, err := api.UpdateRepo(&statistic.UpdateModel{DownloadCount: count}, id)
	if err != nil {
		return xerrors.Errorf("fail to use internal api, error: %w", err)
	}
	return nil
}

type VisitData struct {
	Code int `json:"code"`
	Data []struct {
		Visit  int    `json:"count"`
		RepoID string `json:"repo_id"`
	} `json:"data"`
}

func (d *visit) fetchVisitCounts(url string) (*VisitData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, xerrors.Errorf("fail to fetch data online, error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("fail to read data, error: %w", err)
	}

	var data VisitData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, xerrors.Errorf("fail to Unmarshal data, error: %w", err)
	}
	return &data, nil
}