package util

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// 获取Kube Client 客户端
func KubeClientByConfig(kubeConfigStr string) (*kubernetes.Clientset, *restclient.Config) {

	// 获取kube配置对象
	kubeClientConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigStr)
	if err != nil {
		panic(err.Error())
	}
	// 获取kube客户端
	var kubeClientSet *kubernetes.Clientset
	kubeClientSet, err = kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		panic(err.Error())
	}
	return kubeClientSet, kubeClientConfig
}

// 获取Kube Client 客户端
func KubeClient(cmd *cobra.Command) (*kubernetes.Clientset, *restclient.Config) {

	configStr, _ := cmd.Flags().GetString("kubeconfig")
	return KubeClientByConfig(configStr)
}
