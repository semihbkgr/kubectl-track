package cli

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	objects *[]Object
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
	s := ""
	for _, o := range *m.objects {
		s += o.unstructured.GetName() + "\n"
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("blue")).Render(s)
}

func newModel(objects *[]Object) model {
	return model{objects}
}
