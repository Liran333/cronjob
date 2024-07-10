package moderation

import (
	"github.com/openmerlin/git-access-sdk/filescan"
	userApi "github.com/openmerlin/git-access-sdk/filescan/api"
	"github.com/sirupsen/logrus"
)

const (
	JobModerationReadme = "moderation-readme"
)

func NewModerationReadmeJob(cfg *Config) *Readme {
	return &Readme{cfg: cfg}
}

type Readme struct {
	cfg *Config
}

func (r *Readme) Type() string {
	return JobModerationReadme
}

func (r *Readme) Run() error {
	req := &filescan.ReqToModeration{
		ModerationStatus: "init",
		FileType:         "readme",
		PageNum:          r.cfg.ReadmePageNum,
		UpdatedAt:        r.cfg.ReadmeLastTime,
	}

	logrus.Infof("req: %#v", req)
	_, err := userApi.Cronjob(req)
	if err != nil {
		return err
	}
	return nil
}
