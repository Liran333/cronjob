package main

import (
	"flag"
	"os"

	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/openmerlin/cronjob/config"
	"github.com/openmerlin/cronjob/jobs"
)

const component = "merlin-cronjob"

type options struct {
	service   liboptions.ServiceOptions
	removeCfg bool
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.removeCfg, "rm-cfg", false,
		"whether remove the cfg file after initialized.",
	)

	if err := fs.Parse(args); err != nil {
		fs.PrintDefaults()

		logrus.Fatalf("failed to parse cmdline %s", err)
	}

	return o
}

func main() {
	logrusutil.ComponentInit(component)

	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.Validate(); err != nil {
		logrus.Errorf("Invalid options, err:%s", err.Error())

		return
	}

	// cfg
	cfg, err := config.LoadConfig(o.service.ConfigFile, o.removeCfg)
	if err != nil {
		logrus.Errorf("load config failed, err:%s", err.Error())

		return
	}

	err = jobs.Run(&cfg)
	if err != nil {
		logrus.Errorf("run jobs, err:%s", err.Error())

		return
	}
}
