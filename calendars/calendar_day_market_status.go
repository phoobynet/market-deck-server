package calendars

type CurrentMarketCondition string

const (
	PreMarket          CurrentMarketCondition = "pre_market"
	Open               CurrentMarketCondition = "open"
	PostMarket         CurrentMarketCondition = "post_market"
	ClosedToday        CurrentMarketCondition = "closed_today"
	ClosedForTheDay    CurrentMarketCondition = "closed_for_the_day"
	ClosedOpeningLater CurrentMarketCondition = "closed_opening_later"
)
