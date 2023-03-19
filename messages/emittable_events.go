package messages

type EmittableEvents string

const (
	Snapshots         = "snapshots"
	CalendarDayUpdate = "calendar_day_update"
	Messages          = "messages"
	Errors            = "errors"
)
