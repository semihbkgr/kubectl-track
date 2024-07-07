package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog"
)

func Start(w watch.Interface) error {
	objects := make([]Object, 0)
	m := newModel(&objects)

	go startWatchingEvents(w, &objects)

	klog.Info("starting the program with the model")
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}

func startWatchingEvents(w watch.Interface, objects *[]Object) {
	klog.Infof("starting events watch")
	for {
		e := <-w.ResultChan()
		klog.Infof("event received, type: %s", e.Type)
		obj, err := NewObject(e.Object)
		if err != nil {
			klog.Warning(err)
			continue
		}
		*objects = append(*objects, *obj)
	}
}
