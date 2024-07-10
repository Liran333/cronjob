package moderation

import (
	"github.com/openmerlin/git-access-sdk/filescan"
	userApi "github.com/openmerlin/git-access-sdk/filescan/api"
	"github.com/sirupsen/logrus"
)

const (
	JobModerationInit = "moderation-init"
)

func NewModerationInitJob(cfg *Config) *Moderation {
	return &Moderation{cfg: cfg}
}

type Moderation struct {
	cfg *Config
}

func (d *Moderation) Type() string {
	return JobModerationInit
}

func (d *Moderation) Run() error {
	req := &filescan.ReqToModeration{
		PageNum:   d.cfg.InitPageNum,
		UpdatedAt: d.cfg.InitLastTime,
	}

	logrus.Infof("req: %#v", req)
	_, err := userApi.Cronjob(req)
	if err != nil {
		return err
	}
	return nil
}
