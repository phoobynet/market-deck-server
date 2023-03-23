package messages

type EmittableEvents string

const (
	Snapshots         = "snapshots"
	SnapshotsLite     = "snapshots_lite"
	CalendarDayUpdate = "calendar_day_update"
	Messages          = "messages"
	Errors            = "errors"
)
