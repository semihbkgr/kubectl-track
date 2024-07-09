package cmd

import (
	"flag"
	"io"
	"os"

	track "github.com/semihbkgr/kubectl-track/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	_ "k8s.io/client-go/plugin/pkg/client/auth/exec"
)

var Version string

var rootCmd = &cobra.Command{
	Use:     "kubectl-track <resource> <name>",
	Short:   "capture the changes between resource versions",
	Long:    "kubectl-track monitors and displays changes for a Kubernetes resource, highlighting differences between resource versions to help in debugging and understanding resource evolution",
	Version: Version,
	Args:    cobra.ExactArgs(2),
	RunE:    run,
}

var cf = genericclioptions.NewConfigFlags(false)

func Execute() {
	defer klog.Flush()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().AddGoFlagSet(flag.CommandLine)
	rootCmd.Flags().String("log-file", "", "file to write logs to")

	cf.AddFlags(rootCmd.Flags())
}

func run(cmd *cobra.Command, args []string) error {
	err := configureLogger(cmd.Flags())
	if err != nil {
		return err
	}

	resource := args[0]
	name := args[1]

	klog.Infof("args - resource: %s, name: %s", name, resource)

	opts := track.Options{
		Resource:    resource,
		Name:        name,
		ConfigFlags: cf,
	}

	return opts.Run(cmd.Context())
}

func configureLogger(flagset *pflag.FlagSet) error {
	logFile, err := flagset.GetString("log-file")
	if err != nil {
		return err
	}

	logFlagSet := flag.NewFlagSet("", flag.PanicOnError)
	klog.InitFlags(logFlagSet)
	err = logFlagSet.Set("logtostderr", "false")
	if err != nil {
		return err
	}
	err = logFlagSet.Parse([]string{})
	if err != nil {
		return err
	}

	var w io.Writer = &NopWriter{}
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		w = f
	}
	klog.SetOutput(w)
	return nil
}

type NopWriter struct{}

func (w *NopWriter) Write(p []byte) (n int, err error) {
	return len(p), err
}
