package moderation

import (
	"github.com/openmerlin/git-access-sdk/filescan"
	userApi "github.com/openmerlin/git-access-sdk/filescan/api"
	"github.com/sirupsen/logrus"
)

const (
	JobModerationPic = "moderation-pic"
)

func NewModerationPicJob(cfg *Config) *pic {
	return &pic{cfg: cfg}
}

type pic struct {
	cfg *Config
}

func (p *pic) Type() string {
	return JobModerationPic
}

func (p *pic) Run() error {
	req := &filescan.ReqToModeration{
		ModerationStatus: "init",
		FileType:         "image",
		PageNum:          p.cfg.PicPageNum,
		UpdatedAt:        p.cfg.PicLastTime,
	}

	logrus.Infof("req: %#v", req)
	_, err := userApi.Cronjob(req)
	if err != nil {
		return err
	}
	return nil
}
