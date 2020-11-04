package check

import (
	"fmt"
	"github.com/modood/table"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kube-tools/src/util"
	"strings"
)

func addIngressTLSFlag(flags *pflag.FlagSet) {

}

func NewCmdIngressTLS() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ingresstls",
		Short:                 "检查Ingress 服务使用的TLS证书。",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunIngressTLS(cmd, args)
		},
	}
	addIngressTLSFlag(cmd.Flags())
	return cmd
}

type IngressInfo struct {
	Name      string
	TLS       string
	NameSpace string
}

func RunIngressTLS(cmd *cobra.Command, args []string) {
	kubeClientSet, _ := util.KubeClient(cmd)

	ings, err := kubeClientSet.ExtensionsV1beta1().Ingresses("").List(metav1.ListOptions{})
	if err != nil {
		panic("获取ingress信息失败")
	}
	info := make([]IngressInfo, 0)
	for _, ing := range ings.Items {
		for _, tls := range ing.Spec.TLS {
			ingressInfo := IngressInfo{
				strings.Join(tls.Hosts, ","),
				tls.SecretName,
				ing.Namespace,
			}
			info = append(info, ingressInfo)
		}
	}

	infoStr := table.Table(info)
	fmt.Println(infoStr)
}
