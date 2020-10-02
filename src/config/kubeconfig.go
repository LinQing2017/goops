package config

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

var KubeConfig *restclient.Config
var KubeClientSet *kubernetes.Clientset
var ShellNamespace *string

func InitKube() {

	kubeConfigStr := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "指定config文件")
	ShellNamespace = flag.String("shell-namespaces", "node-shell", "node-shell 的命名空间")
	flag.Parse()

	var err error
	// 获取kube配置对象
	KubeConfig, err = clientcmd.BuildConfigFromFlags("", *kubeConfigStr)
	if err != nil {
		panic(err.Error())
	}
	// 获取kube客户端
	KubeClientSet, err = kubernetes.NewForConfig(KubeConfig)
	if err != nil {
		panic(err.Error())
	}
}
