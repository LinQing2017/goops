package tools

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kube-tools/src/config"
)

func Clean() {

	err := config.KubeClientSet.AppsV1().DaemonSets(*config.ShellNamespace).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{})
	if err != nil {
		fmt.Println("清理命名空间", *config.ShellNamespace, "失败！")
		panic(err.Error())
	}
	err = config.KubeClientSet.CoreV1().Namespaces().Delete(*config.ShellNamespace, &metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("清理命名空间", *config.ShellNamespace, "失败！")
		panic(err.Error())
	}
	fmt.Println("清理命名空间", *config.ShellNamespace, "成功！")
}
