package client

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
)

func NewClient(cf *genericclioptions.ConfigFlags) (*dynamic.DynamicClient, error) {
	config, err := cf.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	return dynamic.NewForConfig(config)
}

func Namespace(cf *genericclioptions.ConfigFlags) string {
	if v := *cf.Namespace; v != "" {
		return v
	}

	clientConfig := cf.ToRawKubeConfigLoader()
	defaultNs, _, err := clientConfig.Namespace()
	if err != nil {
		defaultNs = "default"
	}

	return defaultNs
}
