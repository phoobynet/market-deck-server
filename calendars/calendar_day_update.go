package calendars

type CalendarDayUpdate struct {
	Condition          CurrentMarketCondition `json:"condition"`
	At                 int64                  `json:"at"`
	PreviousMarketDate CalendarDay            `json:"prev"`
	CurrentMarketDate  CalendarDay            `json:"current"`
	NextMarketDate     CalendarDay            `json:"next"`
}
