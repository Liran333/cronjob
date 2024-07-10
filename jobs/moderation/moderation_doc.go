package moderation

import (
	"github.com/openmerlin/git-access-sdk/filescan"
	userApi "github.com/openmerlin/git-access-sdk/filescan/api"
	"github.com/sirupsen/logrus"
)

const (
	JobModerationDoc = "moderation-doc"
)

func NewModerationDocJob(cfg *Config) *Doc {
	return &Doc{cfg: cfg}
}

type Doc struct {
	cfg *Config
}

func (d *Doc) Type() string {
	return JobModerationPic
}

func (d *Doc) Run() error {
	req := &filescan.ReqToModeration{
		ModerationStatus: "init",
		FileType:         "document",
		PageNum:          d.cfg.DocPageNum,
		UpdatedAt:        d.cfg.DocLastTime,
	}

	logrus.Infof("req: %#v", req)
	_, err := userApi.Cronjob(req)
	if err != nil {
		return err
	}
	return nil
}
