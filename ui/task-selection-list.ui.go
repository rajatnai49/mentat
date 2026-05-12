package ui

import (
	"os"
	"os/exec"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/rajatnai49/mentat/vault"
)

type styles struct {
	app           lipgloss.Style
	title         lipgloss.Style
	statusMessage lipgloss.Style
}

func newStyles() styles {
	return styles{
		app: lipgloss.NewStyle().
			Padding(1, 2),
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1),
		statusMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")),
	}
}

type TaskListModel struct {
	styles styles
	list   list.Model
}

func (m TaskListModel) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
	)
}

func (m TaskListModel) View() tea.View {
	v := tea.NewView(m.styles.app.Render(m.list.View()))
	v.AltScreen = true
	return v
}

func (m TaskListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := m.styles.app.GetFrameSize()
		m.list.SetSize(
			msg.Width-h,
			msg.Height-v,
		)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			selected := m.list.SelectedItem().(vault.TaskItem)

			c := exec.Command("nvim", selected.Filepath)
			c.Stdout = os.Stdout
			c.Stdin = os.Stdin
			c.Stderr = os.Stderr

			return m, tea.ExecProcess(c, func(err error) tea.Msg {
				return nil
			})
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func RenderList(tasks []vault.TaskItem) error {
	var items []list.Item

	for _, t := range tasks {
		items = append(items, t)
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Tasks Items"
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)

	p := tea.NewProgram(
		TaskListModel{
			styles: newStyles(),
			list:   l,
		},
	)

	_, err := p.Run()
	return err
}
