package cli

import (
	"bytes"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"k8s.io/cli-runtime/pkg/printers"
)

type model struct {
	resource *Resource
}

func (m model) Init() tea.Cmd {
	return doTick()
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second/24, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case TickMsg:
		return m, doTick()
	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("len: %d\n", len(m.resource.Versions))
	for _, resVersion := range m.resource.Versions {
		s += resVersion.Object.GetName() + "\n"
		printer := printers.NewTablePrinter(printers.PrintOptions{})
		resVersion.Table.ColumnDefinitions = m.resource.TableColumnDefinition()
		buf := bytes.NewBuffer([]byte{})
		printer.PrintObj(resVersion.Table, buf)
		s += buf.String()
		s += "\n------------------------------------------\n"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color("blue")).Render(s)
}

func newModel(resource *Resource) model {
	return model{resource}
}
