package vault

import "time"

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
