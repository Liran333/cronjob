package config

import (
	"os"

	"github.com/openmerlin/cronjob/jobs/downloadcount"
	"github.com/openmerlin/cronjob/utils"
)

func LoadConfig(path string, remove bool) (cfg Config, err error) {
	if remove {
		defer os.Remove(path)
	}

	if err = utils.LoadFromYaml(path, &cfg); err != nil {
		return
	}

	SetDefault(&cfg)

	err = Validate(&cfg)

	return
}

type Config struct {
	DownloadCount downloadcount.Config `json:"download_count"`
}

// ConfigItems returns a slice of interface{} containing pointers to the configuration items.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.DownloadCount,
	}
}

// SetDefault sets default values for the Config struct.
func (cfg *Config) SetDefault() {

}

// Validate validates the configuration.
func (cfg *Config) Validate() error {
	return utils.CheckConfig(cfg, "")
}
