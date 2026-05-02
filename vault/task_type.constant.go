package vault

type TaskType int

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
