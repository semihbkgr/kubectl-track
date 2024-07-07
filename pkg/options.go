package track

import (
	"context"
	"fmt"

	"github.com/semihbkgr/kubectl-track/pkg/cli"
	"github.com/semihbkgr/kubectl-track/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog"
)

type Options struct {
	ConfigFlags *genericclioptions.ConfigFlags
	Resource    string
	Name        string
}

func (o Options) Run(ctx context.Context) error {
	c, err := client.NewClient(o.ConfigFlags)
	if err != nil {
		return fmt.Errorf("error on creating client: %w", err)
	}

	ns := client.Namespace(o.ConfigFlags)

	//todo
	r := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: o.Resource,
	}

	klog.Infof("watching resource in namespace: %s", ns)
	w, err := c.Resource(r).Namespace(ns).Watch(ctx, metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(metav1.ObjectNameField, o.Name).String(),
	})
	if err != nil {
		return fmt.Errorf("error on watching resource: %w", err)
	}

	return cli.Start(w)
}
