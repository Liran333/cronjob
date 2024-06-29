/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package checkrevoke for check user revoke

package checkrevoke

import (
	userApi "github.com/openmerlin/merlin-sdk/user/api"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

const (
	JobCheckUserRevoke = "check—user-revoke"
)

func NewCheckUserRevokeJob() *CheckUserRevoke {
	return &CheckUserRevoke{}
}

type CheckUserRevoke struct {
}

func (c *CheckUserRevoke) Type() string {
	return JobCheckUserRevoke
}

func (c *CheckUserRevoke) Run() error {
	if err := checkUserRevoke(); err != nil {
		logrus.Fatalf("Failed checkUserRevoke: %s", err)
		return err
	}
	return nil
}

func checkUserRevoke() error {
	_, err := userApi.CheckUserRevoke()
	if err != nil {
		return xerrors.Errorf("fail to use internal api, error: %w", err)
	}
	return nil
}
