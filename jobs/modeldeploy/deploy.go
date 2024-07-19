/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package modeldeploy
package modeldeploy

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/client"
	"github.com/opensourceways/server-common-lib/utils"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"

	"github.com/openmerlin/cronjob/common"
	"github.com/openmerlin/merlin-sdk/modeldeploy"
	"github.com/openmerlin/merlin-sdk/modeldeploy/api"
)

const (
	jobExecIntervalHour = 1
	JobModelDeploy      = "deploy-task"
)

type iClient interface {
	GetPathContent(org, repo, path, ref string) (sdk.Content, error)
}

func NewDeployJob(c *Config, gitee common.Gitee) *deploy {
	cli := client.NewClient(func() []byte {
		return []byte(gitee.Token)
	})

	return &deploy{
		cfg:   c,
		cli:   cli,
		gitee: gitee,
	}
}

type deploy struct {
	cfg   *Config
	cli   iClient
	gitee common.Gitee
}

type DeployData struct {
	Cloud  string `json:"cloud"`
	Icon   string `json:"icon"`
	Link   string `json:"link"`
	Desc   string `json:"desc"`
	DescCn string `json:"desc_cn"`
}

func (data *DeployData) SetNewIconUrlOfObs(url string) {
	data.Icon = url
}

func (d deploy) Type() string {
	return JobModelDeploy
}

func (d deploy) Run() error {
	isModified, err := d.isDeployModified()
	if err != nil {
		return err
	}

	if !isModified {
		logrus.Infof("%s is not modified within in %d hour",
			d.cfg.DeployLocation.ToString(), jobExecIntervalHour)

		return nil
	}

	deployContent, err := d.getFileContent(d.cfg.DeployLocation.Path)
	if err != nil {
		return err
	}

	deployData := make(map[string][]DeployData)
	if err = yaml.Unmarshal(deployContent, &deployData); err != nil {
		return err
	}

	for k, data := range deployData {
		modelIndex, err1 := d.getOwnerAndRepo(k)
		if err1 != nil {
			logrus.Error(err1)
			continue
		}

		if err1 = d.handleDeployData(modelIndex, data); err1 != nil {
			logrus.Errorf("handle repo %s failed: %v", modelIndex.ToString(), err1)
		}
	}

	return nil
}

func (d deploy) handleDeployData(index ModelIndex, data []DeployData) error {
	var dataWithNewIconUrl []modeldeploy.RequestToSaveDeploy
	for _, v := range data {
		iconUrl, err := d.generateNewIconUrlOfObs(v)
		if err != nil {
			logrus.Errorf("handle cloud %s in %s error: %v", v.Cloud, index.ToString(), err)
			continue
		}

		v.SetNewIconUrlOfObs(iconUrl)
		modelData := d.toCmd(v)

		dataWithNewIconUrl = append(dataWithNewIconUrl, modelData)
	}

	// todo api to merlin-server
	code, err := api.SaveModelDeploy(index.Owner, index.Repo, dataWithNewIconUrl)
	if err != nil || code != "200" {
		logrus.Errorf("handle cloud in %s error: %v", index.ToString(), err)
	}
	return nil
}

func (d deploy) toCmd(data DeployData) modeldeploy.RequestToSaveDeploy {
	return modeldeploy.RequestToSaveDeploy{
		Cloud:  data.Cloud,
		Icon:   data.Icon,
		Link:   data.Link,
		Desc:   data.Desc,
		DescCn: data.DescCn,
	}
}

func (d deploy) generateNewIconUrlOfObs(v DeployData) (string, error) {
	iconContent, err := d.getFileContent(v.Icon)
	if err != nil {
		return "", fmt.Errorf("get icon content err: %w", err)
	}

	iconUrl, err := d.uploadIconToOBs(iconContent, v.Icon)
	if err != nil {
		return "", fmt.Errorf("upload icon to obs err: %w", err)
	}

	return iconUrl, nil
}

type commitResp struct {
	Commit struct {
		Committer struct {
			Date time.Time `json:"date"`
		} `json:"committer"`
	} `json:"commit"`
}

func (d deploy) isDeployModified() (b bool, err error) {
	commitUrl := d.gitee.CommitUrl(d.cfg.DeployLocation.Owner, d.cfg.DeployLocation.Repo, d.cfg.DeployLocation.Path)

	req, err := http.NewRequest(http.MethodGet, commitUrl, nil)
	if err != nil {
		return
	}

	cli := utils.NewHttpClient(3)
	data, _, err := cli.Download(req)
	if err != nil {
		return
	}

	var resp []commitResp
	if err = json.Unmarshal(data, &resp); err != nil {
		return
	}

	modifiedTime := resp[0].Commit.Committer.Date
	lastOneHour := time.Now().Add(-time.Hour * jobExecIntervalHour)
	b = modifiedTime.After(lastOneHour)

	return
}

func (d deploy) getFileContent(path string) (content []byte, err error) {
	l := d.cfg.DeployLocation
	data, err := d.cli.GetPathContent(l.Owner, l.Repo, path, l.Ref)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(data.Content)
}

func (d deploy) uploadIconToOBs(content []byte, filePath string) (string, error) {
	input := &obs.PutObjectInput{}
	input.Bucket = d.cfg.IconUpload.Bucket
	input.Key = filePath
	input.Body = bytes.NewReader(content)
	out, err := common.ObsClient().PutObject(input)
	if err != nil {
		return "", err
	}
	logrus.Info(out.ObjectUrl)

	split := strings.Split(out.ObjectUrl, "/")
	domain := split[0]
	logrus.Info(domain)
	iconUrl := fmt.Sprintf("%s%s", d.cfg.IconUpload.EndPoint, filePath)

	return iconUrl, nil
}

type ModelIndex struct {
	Owner string
	Repo  string
}

func (m ModelIndex) ToString() string {
	return fmt.Sprintf("[%s/%s]", m.Owner, m.Repo)
}

func (d deploy) getOwnerAndRepo(v string) (ModelIndex, error) {
	split := strings.Split(v, "/")
	if len(split) != 2 {
		return ModelIndex{}, fmt.Errorf("owner/repo format error of %s", v)
	}

	return ModelIndex{
		Owner: split[0],
		Repo:  split[1],
	}, nil
}
