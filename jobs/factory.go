/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package jobs for register factory

package jobs

import (
	"github.com/openmerlin/cronjob/config"
	"github.com/openmerlin/cronjob/jobs/checkrevoke"
	"github.com/openmerlin/cronjob/jobs/downloadcount"
	"github.com/openmerlin/cronjob/jobs/visitcount"
)

type Job interface {
	Run() error
	Type() string
}

var factories = map[string]Job{}

func InitJobMap(cfg *config.Config) {
	factories = map[string]Job{
		downloadcount.JobDownloadCount: downloadcount.NewDownloadJob(&cfg.DownloadCount),
		visitcount.JobVisitCount:       visitcount.NewVisitJob(&cfg.VisitCount),
		checkrevoke.JobCheckUserRevoke: checkrevoke.NewCheckUserRevokeJob(),
	}
}

func GetJobFactory(jobType string) Job {
	var job Job
	job, ok := factories[jobType]
	if !ok {
		return nil
	}
	return job
}
