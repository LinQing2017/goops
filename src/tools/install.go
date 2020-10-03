package tools

import (
	"fmt"
	"github.com/modood/table"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kube-tools/src/config"
	"strings"
	"time"
)

type PodPrintInfo struct {
	Name   string
	Status string
	Node   string
}

func Install() {

	_, err := config.KubeClientSet.AppsV1().DaemonSets(config.ShellNamespace).Get(config.ShellPodName, metav1.GetOptions{
		TypeMeta:        metav1.TypeMeta{},
		ResourceVersion: "",
	})

	if err == nil {
		printPodInfo(config.ShellNamespace, "name="+config.ShellPodName)
		return
	}

	shellNS := v1.Namespace{
		TypeMeta:   metav1.TypeMeta{Kind: "Namespace", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{Name: config.ShellNamespace},
	}
	_, err = config.KubeClientSet.CoreV1().Namespaces().Create(&shellNS)
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
			Name:      config.ShellPodName,
			Namespace: config.ShellNamespace,
		},
		Spec: appv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": config.ShellPodName,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": config.ShellPodName,
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
					},
				},
			},
		},
	}
	_, err = config.KubeClientSet.AppsV1().DaemonSets(config.ShellNamespace).Create(&nodeShellDSDefine)
	if err != nil {
		fmt.Println("创建DaemonSet", nodeShellDSDefine.Name, "失败")
		panic(err.Error())
	}
	time.Sleep(5 * 1e9)
	printPodInfo(config.ShellNamespace, "name="+config.ShellPodName)
}

func printPodInfo(namespace, lableSelector string) {

	pods, err := config.KubeClientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{
		TypeMeta:      metav1.TypeMeta{},
		LabelSelector: lableSelector,
	})
	if err != nil {
		panic(err.Error())
	}
	podInfoList := make([]PodPrintInfo, len(pods.Items))
	for i := 0; i < len(pods.Items); i++ {
		podInfo := PodPrintInfo{
			Name:   pods.Items[i].Name,
			Status: string(pods.Items[i].Status.Phase),
			Node:   pods.Items[i].Status.HostIP,
		}
		podInfoList[i] = podInfo

	}
	table.Output(podInfoList)
}
