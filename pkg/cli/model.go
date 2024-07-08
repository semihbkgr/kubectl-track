package cli

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"k8s.io/cli-runtime/pkg/printers"
)

type model struct {
	resource     *Resource
	cursor       int
	selected     bool
	viewport     viewport.Model
	yamlViewport viewport.Model
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

		if msg.String() == "enter" {
			if !m.selected {
				resVersion := m.resource.Versions[m.cursor]
				y := MapToYaml(resVersion.Object.Object)
				m.yamlViewport.SetContent(y)
				m.yamlViewport.SetYOffset(0)
			}
			m.selected = true
		} else if msg.String() == "esc" {
			m.selected = false
		}

		//todo: scroll the viewport
		if !m.selected {
			if msg.String() == "up" || msg.String() == "k" {
				if m.cursor > 0 {
					m.cursor--
				}
			} else if msg.String() == "down" || msg.String() == "j" {
				if m.cursor < len(m.resource.Versions)-1 {
					m.cursor++
				}
			}
		} else {
			var cmd tea.Cmd
			m.yamlViewport, cmd = m.yamlViewport.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		m.yamlViewport.Width = msg.Width
		m.yamlViewport.Height = msg.Height
	case TickMsg:
		return m, doTick()
	}

	return m, nil
}

func (m model) View() string {
	if m.selected {
		return m.yamlViewport.View()
	}

	b := strings.Builder{}
	for i, resVersion := range m.resource.Versions {
		var bgColor lipgloss.TerminalColor = lipgloss.NoColor{}
		if i == m.cursor {
			bgColor = lipgloss.ANSIColor(239)
		}

		title := fmt.Sprintf("%s - %s", resVersion.Version, resVersion.Timestamp.Format(time.RFC3339))
		b.Write([]byte(lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(85)).Background(bgColor).Render(title)))
		if resVersion.Table != nil {
			b.WriteByte('\n')
			//todo: colorize table
			printer := printers.NewTablePrinter(printers.PrintOptions{})
			resVersion.Table.ColumnDefinitions = m.resource.TableColumnDefinition()
			buf := bytes.NewBuffer([]byte{})
			printer.PrintObj(resVersion.Table, buf)
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(153)).Render(buf.String()))
		}
		b.WriteByte('\n')
	}

	m.viewport.SetContent(b.String())

	return m.viewport.View()
}

func newModel(resource *Resource) model {
	return model{
		resource:     resource,
		viewport:     viewport.New(0, 0),
		yamlViewport: viewport.New(0, 0),
	}
}
