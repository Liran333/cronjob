package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"sigs.k8s.io/yaml"

	"github.com/openmerlin/merlin-sdk/httpclient"
	"github.com/openmerlin/merlin-sdk/statistic"
	"github.com/openmerlin/merlin-sdk/statistic/api"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
)

const component = "merlin-cronjob"

type Options struct {
	Service   liboptions.ServiceOptions
	RemoveCfg bool
}

func (o *Options) AddFlags(fs *flag.FlagSet) {
	o.Service.AddFlags(fs)
	fs.BoolVar(&o.RemoveCfg, "rm-cfg", false, "Remove the cfg file after initialization.")
}

type SDKConfig = httpclient.Config

type Config struct {
	Merlin        SDKConfig      `json:"merlin"`
	DownloadCount DownloadConfig `json:"download_count"`
	VisitCount    VisitConfig    `json:"visit_count"`
}

type DownloadConfig struct {
	Spec            string `json:"spec" required:"true"`
	OriginalDataUrl string `json:"original_data_url" required:"true"`
}

type VisitConfig struct {
	Spec            string `json:"spec" required:"true"`
	OriginalDataUrl string `json:"original_data_url" required:"true"`
}

func LoadConfig(path string, remove bool) (*Config, error) {
	var cfg Config
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	if remove {
		if err := os.Remove(path); err != nil {
			return nil, err
		}
	}
	return &cfg, nil
}

type DownloadData struct {
	Code int `json:"code"`
	Data []struct {
		Name     string `json:"name"`
		Download int    `json:"download"`
		RepoID   string `json:"repo_id"`
	} `json:"data"`
}

type VisitData struct {
	Code int `json:"code"`
	Data []struct {
		Visit  int    `json:"count"`
		RepoID string `json:"repo_id"`
	} `json:"data"`
}

func fetchDownloadCounts(url string) (*DownloadData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, xerrors.Errorf("fail to fetch data online, error: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("fail to read data, error: %w", err)
	}

	var data DownloadData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, xerrors.Errorf("fail to Unmarshal data, error: %w", err)
	}
	return &data, nil
}

func fetchVisitCounts(url string) (*VisitData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, xerrors.Errorf("fail to fetch data online, error: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("fail to read data, error: %w", err)
	}

	var data VisitData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, xerrors.Errorf("fail to Unmarshal data, error: %w", err)
	}
	return &data, nil
}

func updateRepo(id string, data *statistic.UpdateModel) error {
	_, err := api.UpdateRepo(data, id)
	if err != nil {
		return xerrors.Errorf("fail to use internal api, error: %w", err)
	}
	return nil
}

func buildAndUpdateRepo(cfg *Config) {
	downloadData, err := fetchDownloadCounts(cfg.DownloadCount.OriginalDataUrl)
	if err != nil {
		logrus.Errorf("Error fetching download counts: %s", err)
		return
	}
	visitData, err := fetchVisitCounts(cfg.VisitCount.OriginalDataUrl)
	if err != nil {
		logrus.Errorf("Error fetching visit counts: %s", err)
	}

	repoData := make(map[string]statistic.UpdateModel)

	for _, repo := range downloadData.Data {
		repoData[repo.RepoID] = statistic.UpdateModel{
			DownloadCount: repo.Download,
			VisitCount:    0,
		}
	}

	for _, repo := range visitData.Data {
		cur, ok := repoData[repo.RepoID]
		if ok {
			cur.VisitCount = repo.Visit
		} else {
			repoData[repo.RepoID] = statistic.UpdateModel{
				DownloadCount: 0,
				VisitCount:    repo.Visit,
			}
		}
	}

	for id, data := range repoData {
		if err := updateRepo(id, &data); err != nil {
			logrus.Errorf("Error update data: %s", err)
		}
	}
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

	cfg, err := LoadConfig(opts.Service.ConfigFile, opts.RemoveCfg)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %s", err)
		return
	}

	httpclient.Init(&cfg.Merlin)

	buildAndUpdateRepo(cfg)

}
