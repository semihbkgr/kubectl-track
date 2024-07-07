package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog"
)

func Start(watchRes watch.Interface, watchResTable watch.Interface) error {
	resource := NewResource()
	m := newModel(resource)

	go startWatchingResourceEvents(resource, watchRes, watchResTable)

	klog.Info("start the program with the model")
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}

func startWatchingResourceEvents(resource *Resource, watchObject watch.Interface, watchTable watch.Interface) {
	klog.Infof("start watching resource events")
	objChan := watchObject.ResultChan()
	tableChan := watchTable.ResultChan()
	for {
		select {
		case objEvent := <-objChan:
			klog.Infof("object event received, type: %s", objEvent.Type)
			obj := objEvent.Object
			if obj == nil {
				klog.Info("object in the event is nil")
				continue
			}
			unsObj := obj.(*unstructured.Unstructured)
			version := unsObj.GetResourceVersion()
			klog.Infof("resource version: %s", version)
			resVersion := resource.CreateOrGetVersion(version)
			resVersion.Object = unsObj
		case tableEvent := <-tableChan:
			klog.Infof("table event received, type: %s", tableEvent.Type)
			obj := tableEvent.Object
			if obj == nil {
				klog.Info("object in the event is nil")
				continue
			}
			table, err := DecodeIntoTable(obj)
			if err != nil {
				panic(err)
			}
			version := table.GetResourceVersion()
			//version := table.Rows[0].Object.Object.(*unstructured.Unstructured).GetResourceVersion()
			klog.Infof("resource version: %s", version)
			resVersion := resource.CreateOrGetVersion(version)
			resVersion.Table = table
		}
	}
}
