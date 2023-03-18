package events

type EmittableEvents string

const (
	RealtimeSymbols   = "realtime_symbols"
	CalendarDayUpdate = "calendar_day_update"
	Messages          = "messages"
	Errors            = "errors"
)
