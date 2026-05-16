package ui

import (
	"fmt"
	"os"
	"os/exec"

	"charm.land/bubbles/v2/key"
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

type editorErrMsg struct {
	err error
}

func newStyles() styles {
	return styles{
		app: lipgloss.NewStyle().
			Padding(1, 2),
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f4ede8")).
			Background(lipgloss.Color("#286983")).
			Padding(0, 1),
		statusMessage: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ea9d34")),
	}
}

type TaskListModel struct {
	styles styles
	list   list.Model
	loadFn func() ([]vault.TaskItem, error)
	cfg *vault.Config
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
		case "enter":
			if m.list.FilterState() == list.Filtering {
				break
			}
			selected := m.list.SelectedItem().(vault.TaskItem)

			c := exec.Command(m.cfg.Editor, selected.Filepath)
			c.Stdout = os.Stdout
			c.Stdin = os.Stdin
			c.Stderr = os.Stderr

			return m, tea.ExecProcess(c, func(err error) tea.Msg {
				if err != nil {
					return editorErrMsg{err: err}
				}
				return nil
			})
		case "r":
			if m.list.FilterState() == list.Filtering {
				break
			}
			m.refresh()
		}

	case editorErrMsg:
		m.list.NewStatusMessage(
			m.styles.statusMessage.Render(
				fmt.Sprintf("failed opening editor: %v", msg.err),
			),
		)
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *TaskListModel) refresh() {
	tasks, _ := m.loadFn()
	var items []list.Item

	for _, t := range tasks {
		items = append(items, t)
	}

	m.list.SetItems(items)
}

func RenderList(cfg *vault.Config, loadFn func() ([]vault.TaskItem, error)) error {
	tasks, err := loadFn()
	if err != nil {
		return err
	}

	var items []list.Item

	for _, t := range tasks {
		items = append(items, t)
	}

	d := getMentatDelegate()

	l := list.New(items, d, 0, 0)
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f4ede8")).
		Background(lipgloss.Color("#eb6f92")).
		Padding(0, 2).
		Bold(true)

	l.Title = "Task Items"
	l.Styles.Title = titleStyle
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "refresh"),
			),
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "open"),
			),
		}
	}

	p := tea.NewProgram(
		TaskListModel{
			styles: newStyles(),
			list:   l,
			loadFn: loadFn,
			cfg: cfg,
		},
	)

	_, err = p.Run()
	return err
}

func getMentatDelegate() list.ItemDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#ebbcba")).
		BorderForeground(lipgloss.Color("#c4a7e7")).
		Bold(true)

	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#908caa")).
		BorderForeground(lipgloss.Color("#c4a7e7"))

	d.Styles.NormalTitle = d.Styles.NormalTitle.
		Foreground(lipgloss.Color("#e0def4"))

	d.Styles.NormalDesc = d.Styles.NormalDesc.
		Foreground(lipgloss.Color("#6e6a86"))

	d.Styles.DimmedTitle = d.Styles.DimmedTitle.
		Foreground(lipgloss.Color("#6e6a86"))

	return d
}
