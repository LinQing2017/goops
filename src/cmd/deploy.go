package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kube-tools/src/config"
	"kube-tools/src/util"
	"strings"
	"time"
)

func RunInstall(cmd *cobra.Command, args []string) {

	kubeClientSet, _ := util.KubeClient(cmd)

	_, err := kubeClientSet.AppsV1().DaemonSets(config.ShellNamespace).Get(config.ShellDaemonset, metav1.GetOptions{
		TypeMeta:        metav1.TypeMeta{},
		ResourceVersion: "",
	})

	if err == nil {
		util.PrintPodSimpleInfo(kubeClientSet, config.ShellNamespace, "name="+config.ShellDaemonset)
		return
	}

	shellNS := v1.Namespace{
		TypeMeta:   metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: config.ShellNamespace},
	}
	_, err = kubeClientSet.CoreV1().Namespaces().Create(&shellNS)
	if err != nil {
		fmt.Println("创建命名空间", config.ShellNamespace, "失败")
		panic(err.Error())
	}

	isPrivileged := true
	var priority int32 = 0
	nodeShellDSDefine := appv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "apps/v1",
			APIVersion: "DaemonSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.ShellDaemonset,
			Namespace: config.ShellNamespace,
		},
		Spec: appv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": config.ShellDaemonset,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": config.ShellDaemonset,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            "shell",
							Image:           "docker.io/alpine:3.9",
							ImagePullPolicy: "IfNotPresent",
							Command:         strings.Fields("nsenter"),
							Args:            strings.Fields("-t 1 -m -u -i -n sleep inf"),
							SecurityContext: &v1.SecurityContext{
								Privileged: &isPrivileged,
							},
							WorkingDir: "/root",
						},
					},
					DNSPolicy:     v1.DNSClusterFirst,
					HostIPC:       true,
					HostNetwork:   true,
					HostPID:       true,
					Priority:      &priority,
					RestartPolicy: v1.RestartPolicyAlways,
					Tolerations: []v1.Toleration{
						{
							Key:      "node-role.kubernetes.io/master",
							Operator: v1.TolerationOpExists,
							Effect:   v1.TaintEffectNoSchedule,
						},
						{
							Key:      "build",
							Value:    "type",
							Operator: v1.TolerationOpEqual,
							Effect:   v1.TaintEffectNoExecute,
						},
					},
				},
			},
		},
	}
	_, err = kubeClientSet.AppsV1().DaemonSets(config.ShellNamespace).Create(&nodeShellDSDefine)
	if err != nil {
		fmt.Println("创建DaemonSet", nodeShellDSDefine.Name, "失败")
		panic(err.Error())
	}
	time.Sleep(5 * 1e9)
	util.PrintPodSimpleInfo(kubeClientSet, config.ShellNamespace, "name="+config.ShellDaemonset)
}

func RunClean(cmd *cobra.Command, args []string) {

	kubeClientSet, _ := util.KubeClient(cmd)

	_, err := kubeClientSet.CoreV1().Namespaces().Get(config.ShellNamespace, metav1.GetOptions{
		TypeMeta: metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
	})
	if err != nil {
		fmt.Println("清理命名空间", config.ShellNamespace, "成功！")
		return
	}
	err = kubeClientSet.AppsV1().DaemonSets(config.ShellNamespace).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{})
	if err != nil {
		fmt.Println("清理命名空间", config.ShellNamespace, "失败！")
		panic(err.Error())
	}
	err = kubeClientSet.CoreV1().Namespaces().Delete(config.ShellNamespace, &metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("清理命名空间", config.ShellNamespace, "失败！")
		panic(err.Error())
	}
	fmt.Println("清理命名空间", config.ShellNamespace, "成功！")
}

func NewCmdInstall() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "install",
		Short:                 "安装Shell Pod服务",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunInstall(cmd, args)
		},
	}
	return cmd
}

func NewCmdClean() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "clean",
		Short:                 "清理Shell Pod服务",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunClean(cmd, args)
		},
	}
	return cmd
}

func NewCmdDeploy() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "deploy [command]",
		Short:                 "在集群部署基于Daemonset的Shell客户端",
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(NewCmdInstall())
	cmd.AddCommand(NewCmdClean())
	return cmd
}
