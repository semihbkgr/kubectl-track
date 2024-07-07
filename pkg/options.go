package track

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/semihbkgr/kubectl-track/pkg/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/rest"
)

type Options struct {
	ConfigFlags *genericclioptions.ConfigFlags
	Resource    string
	Name        string
}

func (o Options) Run(ctx context.Context) error {
	res, err := o.resource()
	if err != nil {
		return fmt.Errorf("error on getting resource: %w", err)
	}
	resTable, err := o.resourceInTable()
	if err != nil {
		return fmt.Errorf("error on getting resource table: %w", err)
	}

	watchRes, err := res.Watch("0")
	if err != nil {
		return fmt.Errorf("error on watching resource: %w", err)
	}
	watchResTable, err := resTable.Watch("0")
	if err != nil {
		return fmt.Errorf("error on watching resource table: %w", err)
	}

	return cli.Start(watchRes, watchResTable)
}

func (o Options) resource() (*resource.Result, error) {
	ns, err := o.namespace()
	if err != nil {
		return nil, err
	}

	return resource.NewBuilder(o.ConfigFlags).
		Unstructured().
		NamespaceParam(ns).
		DefaultNamespace().
		ResourceNames(o.Resource, o.Name).
		SingleResourceType().
		Latest().
		Do(), nil
}

func (o Options) resourceInTable() (*resource.Result, error) {
	ns, err := o.namespace()
	if err != nil {
		return nil, err
	}

	return resource.NewBuilder(o.ConfigFlags).
		Unstructured().
		NamespaceParam(ns).
		DefaultNamespace().
		ResourceNames(o.Resource, o.Name).
		SingleResourceType().
		Latest().
		TransformRequests(transformRequestsAcceptTableHeader).
		Do(), nil
}

func transformRequestsAcceptTableHeader(req *rest.Request) {
	req.SetHeader("Accept", strings.Join([]string{
		fmt.Sprintf("application/json;as=Table;v=%s;g=%s", metav1.SchemeGroupVersion.Version, metav1.GroupName),
	}, ","))
}

func (o Options) namespace() (string, error) {
	ns, _, err := o.ConfigFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return "", err
	}
	if ns == "" {
		return "", errors.New("empty namespace")
	}
	return ns, nil
}
