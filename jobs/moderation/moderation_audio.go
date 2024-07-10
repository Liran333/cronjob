package moderation

import (
	"github.com/openmerlin/git-access-sdk/filescan"
	userApi "github.com/openmerlin/git-access-sdk/filescan/api"
	"github.com/sirupsen/logrus"
)

const (
	JobModerationAudio = "moderation-audio"
)

func NewModerationAudioJob(cfg *Config) *audio {
	return &audio{cfg: cfg}
}

type audio struct {
	cfg *Config
}

func (a *audio) Type() string {
	return JobModerationPic
}

func (a *audio) Run() error {
	req := &filescan.ReqToModeration{
		ModerationStatus: "init",
		FileType:         "audio",
		PageNum:          a.cfg.AudioPageNum,
		UpdatedAt:        a.cfg.AudioLastTime,
	}

	logrus.Infof("req: %#v", req)
	_, err := userApi.Cronjob(req)
	if err != nil {
		return err
	}
	return nil
}
