package jobs

import (
	"github.com/robfig/cron/v3"

	"github.com/openmerlin/cronjob/config"
)

type Job interface {
	Spec() string
	Handle()
}

func Run(cfg *config.Config) error {

}

func run(jobs []Job) error {
	c := cron.New()

	for _, job := range jobs {
		_, err := c.AddFunc(job.Spec(), job.Handle)
		if err != nil {

		}
	}

	c.Start()

	return nil
}