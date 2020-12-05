package ews_client

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func GetCluster(appname, clusterid string) (EWSCluster, error) {

	cluster := EWSCluster{}

	url := "http://" + EWSURL + "/v1/apps/" + appname + "/clusters/" + clusterid
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorf("构造Request失败（%s）。", err.Error())
		return cluster, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+EWSToken)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		logrus.Errorf("EWS服务（%s，%s）请求失败。", appname, clusterid, err)
		return cluster, err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(rbody, &cluster)
	return cluster, err
}
