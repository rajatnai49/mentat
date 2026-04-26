package vault

type TaskStatus int
type TaskType int

const (
	Open TaskStatus = iota
	Done
	OnHold
)

var TaskStatusName = map[TaskStatus]string{
	Open:   "open",
	Done:   "done",
	OnHold: "on-hold",
}

const (
	StandAlone TaskType = iota
	Daily
	Monthly
	Yearly
)

var TaskTypeName = map[TaskType]string{
	StandAlone: "general",
	Daily: "daily",
	Monthly: "monthly",
	Yearly: "yearly",
}
