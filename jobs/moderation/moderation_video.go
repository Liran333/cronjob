package moderation

import (
	"github.com/openmerlin/git-access-sdk/filescan"
	userApi "github.com/openmerlin/git-access-sdk/filescan/api"
	"github.com/sirupsen/logrus"
)

const (
	JobModerationVideo = "moderation-video"
)

func NewModerationVideoJob(cfg *Config) *Video {
	return &Video{cfg: cfg}
}

type Video struct {
	cfg *Config
}

func (v *Video) Type() string {
	return JobModerationPic
}

func (v *Video) Run() error {
	req := &filescan.ReqToModeration{
		ModerationStatus: "init",
		FileType:         "video",
		PageNum:          v.cfg.VideoPageNum,
		UpdatedAt:        v.cfg.VideoLastTime,
	}

	logrus.Infof("req: %#v", req)
	_, err := userApi.Cronjob(req)
	if err != nil {
		return err
	}
	return nil
}
