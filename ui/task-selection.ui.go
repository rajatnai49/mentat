package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/rajatnai49/mentat/vault"
)

type TaskUIModel struct {
	cursor int
	tasks  []vault.TaskItem
	choice vault.TaskItem
}

func (tm TaskUIModel) Init() tea.Cmd {
	return nil
}

func (tm TaskUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return tm, tea.Quit
		case "enter":
			selected := tm.tasks[tm.cursor]

			cmd := exec.Command("nvim", selected.Filepath)
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr

			return tm, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return nil
			})
		case "down", "j":
			if tm.cursor < len(tm.tasks)-1 {
				tm.cursor++
			}
		case "up", "k":
			if tm.cursor > 0 {
				tm.cursor--
			}
		}
	}
	return tm, nil
}

func (tm TaskUIModel) View() tea.View {
	s := strings.Builder{}
	s.WriteString("\nSelect Task you want to open \n\n")

	for i, task := range tm.tasks {
		if tm.cursor == i {
			s.WriteString("> ")
		}
		s.WriteString(fmt.Sprintf("%v. %s\n", i+1, task.Task.Title))
		s.WriteString("\n")
	}
	s.WriteString("\n[j/k] move • [enter] select • [q] quit")

	v := tea.NewView(s.String())
	v.AltScreen = true

	return v
}

func RenderTasks(tasks []vault.TaskItem) error {
	p := tea.NewProgram(
		TaskUIModel{
			tasks: tasks,
		},
	)

	_, err := p.Run()
	if err != nil {
		return err
	}

	// 	_, ok := m.(TaskUIModel)
	// 	if !ok {
	// 		return fmt.Errorf("Invalid model")
	// 	}

	return nil
}
