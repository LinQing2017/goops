package util

import (
	"context"
	"github.com/docker/docker/client"
	"path"
	"strings"
)

var (
	defaultAPIVerison = "1.25"
)

// 返回Docker Client
func DockerClient(host string) *client.Client {

	var c *client.Client
	var err error
	if strings.EqualFold(host, "") {
		c, err = client.NewEnvClient()
	} else {
		c, err = client.NewClient(host, defaultAPIVerison, nil, nil)
	}
	if err != nil {
		panic(err)
	}

	return c
}

// 获取容器运行时占用的磁盘空间
func ContainerSize(containerID string, cli *client.Client) int64 {
	containerInfo, _ := cli.ContainerInspect(context.Background(), containerID)
	upperDir := containerInfo.GraphDriver.Data["UpperDir"]
	upperDirSize, _ := CalculateDirSize(upperDir)
	return upperDirSize
}

// 获取容器日志的磁盘使用空间
func ContainerLogSize(containerID string, cli *client.Client) int64 {
	dockerInfo, _ := cli.Info(context.Background())
	containerDataPath := path.Join(dockerInfo.DockerRootDir, containerID)
	logSize, _ := CalculateDirSize(containerDataPath)
	return logSize
}
