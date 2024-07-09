package cli

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/printers"
)

type model struct {
	resource     *Resource
	cursor       int
	selected     bool
	viewport     viewport.Model
	yamlViewport viewport.Model
	rvTableCache map[string]string
	rvDiffCache  map[string]string
	rvYamlCache  map[string]string
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
				if _, ok := m.rvYamlCache[resVersion.Version]; !ok {
					truncated := truncateObjectForYaml(*resVersion.Object)
					y := YamlRenderString(truncated.Object)
					m.rvYamlCache[resVersion.Version] = y
				}
				m.yamlViewport.SetContent(m.rvYamlCache[resVersion.Version])
				m.yamlViewport.SetYOffset(0)
			}
			m.selected = true
		} else if msg.String() == "esc" {
			m.selected = false
		}

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

		title := fmt.Sprintf("● %s - %s", resVersion.Version, resVersion.Timestamp.Format(time.DateTime))
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(85)).Background(bgColor).Render(title))
		b.WriteByte('\n')
		if resVersion.Table != nil {
			//todo: colorize table
			if _, ok := m.rvTableCache[resVersion.Version]; !ok {
				printer := printers.NewTablePrinter(printers.PrintOptions{Wide: true})
				resVersion.Table.ColumnDefinitions = m.resource.TableColumnDefinition()
				buf := bytes.NewBuffer([]byte{})
				printer.PrintObj(resVersion.Table, buf)
				t := lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.ANSIColor(153)).Render(buf.String())
				m.rvTableCache[resVersion.Version] = t
			}

			b.WriteString(m.rvTableCache[resVersion.Version])
		}

		if i == m.cursor && i > 0 && resVersion.Object != nil {
			if m.resource.Versions[i-1].Object.GetUID() == resVersion.Object.GetUID() {
				b.WriteByte('\n')
				if _, ok := m.rvDiffCache[resVersion.Version]; !ok {
					diff := DiffRenderString(truncateObjectForDiff(*m.resource.Versions[i-1].Object).Object, truncateObjectForDiff(*resVersion.Object).Object)
					d := lipgloss.NewStyle().PaddingLeft(2).Render(diff)
					m.rvDiffCache[resVersion.Version] = d
				}
				b.WriteString(m.rvDiffCache[resVersion.Version])
			}
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
		rvTableCache: make(map[string]string),
		rvDiffCache:  make(map[string]string),
		rvYamlCache:  make(map[string]string),
	}
}

func truncateObjectForYaml(u unstructured.Unstructured) unstructured.Unstructured {
	u.SetManagedFields(nil)
	return u
}

func truncateObjectForDiff(u unstructured.Unstructured) unstructured.Unstructured {
	u.SetManagedFields(nil)
	u.SetResourceVersion("")
	return u
}
