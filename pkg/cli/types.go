package cli

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type Resource struct {
	VersionsMap map[string]*ResourceVersion
	Versions    []*ResourceVersion
}

func NewResource() *Resource {
	return &Resource{
		VersionsMap: make(map[string]*ResourceVersion),
		Versions:    make([]*ResourceVersion, 0),
	}
}

func (r *Resource) CreateOrGetVersion(version string) *ResourceVersion {
	if v, ok := r.VersionsMap[version]; ok {
		return v
	}
	v := NewResourceVersion(version, time.Now())
	r.VersionsMap[version] = v
	r.Versions = append(r.Versions, v)
	return v
}

func (r *Resource) TableColumnDefinition() []metav1.TableColumnDefinition {
	return r.Versions[0].Table.ColumnDefinitions
}

type ResourceVersion struct {
	Version   string
	Timestamp time.Time
	EventType watch.EventType
	Object    *unstructured.Unstructured
	Table     *metav1.Table
}

func NewResourceVersion(version string, timestamp time.Time) *ResourceVersion {
	//todo: timestamp
	return &ResourceVersion{
		Version:   version,
		Timestamp: timestamp,
	}
}

func DecodeIntoTable(o runtime.Object) (*metav1.Table, error) {
	if o.GetObjectKind().GroupVersionKind() != metav1.SchemeGroupVersion.WithKind("Table") {
		return nil, fmt.Errorf("cannot decode non-table object into table")
	}

	u, ok := o.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("cannot decode non-unstructured object into table")
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
