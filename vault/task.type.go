package vault

import (
	"strings"
	"time"
)

type Task struct {
	Title       string
	Description string
	Done        bool
	Type        TaskType
	Tags        []string
	LinkedNotes []string
}

type NoteTask struct {
	Date     time.Time
	FilePath string
	Tasks    []Task
}

type TaskItem struct {
	Task     Task
	Filepath string
}

func (i TaskItem) Title() string {
	return i.Task.Title
}

func (i TaskItem) Description() string {
	var parts []string

	if len(i.Task.Tags) > 0 {
		parts = append(
			parts,
			"#"+strings.Join(i.Task.Tags, " #"),
		)
	}

	if i.Task.Description != "" {
		parts = append(parts, i.Task.Description)
	}

	return strings.Join(parts, " • ")
}

func (i TaskItem) GetFilePath() string {
	return i.Filepath
}

func (i TaskItem) FilterValue() string {
	t := i.Task

	return t.Title + " " + strings.Join(t.Tags, " ")
}
