package vault

import "time"

type Task struct {
	ID          string
	Title       string
	Description string
	Type        TaskType
	Status      TaskStatus
	HoldReason  string
	Depends     []string
	Parent      string
	Tags        []string
}

type NoteTask struct {
	Date     time.Time
	FilePath string
	Tasks    []Task
}

func (t Task) IsOnHold() bool {
	return t.Status == OnHold
}
