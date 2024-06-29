/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package main for cron job

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/openmerlin/merlin-sdk/httpclient"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/cronjob/config"
	"github.com/openmerlin/cronjob/jobs"
)

const component = "merlin-cronjob"

type Options struct {
	Service   liboptions.ServiceOptions
	JobType   string
	RemoveCfg bool
}

func (o *Options) Validate() error {
	if o.JobType == "" {
		return fmt.Errorf("missing job type")
	}
	return o.Service.Validate()
}

func (o *Options) AddFlags(fs *flag.FlagSet) {
	o.Service.AddFlags(fs)
	fs.BoolVar(&o.RemoveCfg, "rm-cfg", false, "Remove the cfg file after initialization.")
	fs.StringVar(&o.JobType, "job-type", "", "type of job")
}

func main() {
	logrusutil.ComponentInit(component)

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	opts := &Options{}
	opts.AddFlags(fs)
	if err := fs.Parse(os.Args[1:]); err != nil {
		logrus.Fatalf("Failed to parse command line: %s", err)
		return
	}

	if err := opts.Validate(); err != nil {
		logrus.Fatalf("failed to validate: %v", err)

		return
	}

	cfg := new(config.Config)
	err := config.LoadConfig(opts.Service.ConfigFile, cfg, opts.RemoveCfg)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %s", err)
		return
	}

	httpclient.Init(&cfg.Merlin)
	jobs.InitJobMap(cfg)

	job := jobs.GetJobFactory(opts.JobType)
	if job == nil {
		logrus.Info("no find job to exec")
		return
	}
	if err = job.Run(); err != nil {
		logrus.Errorf("run job %s failed: %s", job.Type(), err.Error())
	}
	return
}
