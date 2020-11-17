package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

// 获取Kube Client 客户端
func KubeClientByConfig(configStr string) *kubernetes.Clientset {
	var err error
	var config *restclient.Config
	if strings.EqualFold(configStr, "") {
		config, err = restclient.InClusterConfig()
	} else {
		// 获取kube配置对象
		config, err = clientcmd.BuildConfigFromFlags("", configStr)
	}
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset

}

// 获取Kube Client 客户端
func KubeClientAndConfig(configStr string) (*kubernetes.Clientset, *restclient.Config) {
	kubeClientConfig, err := clientcmd.BuildConfigFromFlags("", configStr)
	if err != nil {
		panic(err.Error())
	}
	kubeClientSet, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		panic(err.Error())
	}
	return kubeClientSet, kubeClientConfig
}
