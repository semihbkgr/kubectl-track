package cli

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type Object struct {
	unstructured *unstructured.Unstructured
	table        *metav1.Table
}

func NewObject(o runtime.Object) (*Object, error) {
	if o == nil {
		return &Object{}, fmt.Errorf("runtime object is nil")
	}

	u, ok := o.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("runtime object is not unstructured")
	}

	table, err := decodeIntoTable(o)
	if err != nil {
		return nil, fmt.Errorf("error in decode table: %v", err)
	}

	return &Object{
		unstructured: u,
		table:        table,
	}, nil
}

func decodeIntoTable(o runtime.Object) (*metav1.Table, error) {
	u, ok := o.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("attempt to decode non-unstructured object")
	}

	table := &metav1.Table{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, table); err != nil {
		return nil, err
	}

	for i := range table.Rows {
		row := &table.Rows[i]
		if row.Object.Raw == nil || row.Object.Object != nil {
			continue
		}
		converted, err := runtime.Decode(unstructured.UnstructuredJSONScheme, row.Object.Raw)
		if err != nil {
			return nil, err
		}
		row.Object.Object = converted
	}

	return table, nil
}
