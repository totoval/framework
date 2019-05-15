package schedule

type When = string

const (
	EveryYear   When = "@yearly"  // Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
	EveryMonth       = "@monthly" // Run once a month, midnight, first of month | 0 0 0 1 * *
	EveryWeek        = "@weekly"  // Run once a week, midnight between Sat/Sun  | 0 0 0 * * 0
	EveryDay         = "@daily"   // Run once a day, midnight                   | 0 0 0 * * *
	EveryHour        = "@hourly"  // Run once an hour, beginning of hour        | 0 0 * * * *
	Every            = "@every"
	EveryMinute      = Every + " " + "1m" // Run once a minute, beginning of minute        | 0 0 * * * *
	EverySecond      = Every + " " + "1s" // Run once a second, beginning of second

	EveryHoursFormat   = Every + " " + "%dh"
	EveryMinutesFormat = Every + " " + "%dm"
	EverySecondsFormat = Every + " " + "%ds"
)
