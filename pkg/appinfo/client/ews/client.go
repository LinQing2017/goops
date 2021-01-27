package ews

import (
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func GetCluster(appname, clusterid string) (EWSCluster, error) {

	var cluster EWSCluster

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+EWSToken).
		ForceContentType("application/json").
		SetResult(&cluster).
		Get("http://" + EWSURL + "/v1/apps/" + appname + "/clusters/" + clusterid)

	if err != nil || resp.StatusCode() != 200 {
		logrus.Error(errors.Wrapf(err, "EWS服务（%s，%s）请求失败。", appname, clusterid, err).Error())
		return cluster, err
	}
	return cluster, nil
}
