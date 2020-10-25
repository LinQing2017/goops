package config

import (
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

var (
	kubeConfigStr  string
	ShellNamespace string
	ShellDaemonset string
)

func AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&kubeConfigStr, "kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Path to the kubeconfig file to use for CLI requests.")

	flags.StringVar(&ShellNamespace, "shellnamespace", "node-shell", "Shell pod namespaces.")
	flags.StringVar(&ShellDaemonset, "shelldeploy", "node-shell-tool", "Shell pod daemonset.")
}
