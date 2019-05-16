package zone

import (
	"time"

	"github.com/totoval/framework/config"
)

var location *time.Location

type Time = time.Time
type Duration = time.Duration
type Location = time.Location
type Timer = time.Timer
type Ticker = time.Ticker

type Month = time.Month
type Weekday = time.Weekday

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)

func Initialize() {
	var err error
	location, err = time.LoadLocation(config.GetString("app.timezone"))
	if err != nil {
		panic(err)
	}
}
func Date(year int, month Month, day, hour, min, sec, nsec int) Time {
	return time.Date(year, month, day, hour, min, sec, nsec, GetLocation())
}
func Until(t Time) Duration {
	return time.Until(t.In(GetLocation()))
}
func AfterFunc(d Duration, f func()) *Timer {
	return time.AfterFunc(d, f)
}
func After(d Duration) <-chan Time {
	t := <-time.After(d)
	t.In(GetLocation())
	r := make(chan Time, 1)
	r <- t
	return r
}
func Since(t Time) Duration {
	return time.Since(t)
}
func Now() Time {
	return time.Now().In(GetLocation())
}
func NewTimer(d Duration) *Timer {
	return time.NewTimer(d)
}
func NewTicker(d Duration) *Ticker {
	return time.NewTicker(d)
}
func Tick(d Duration) <-chan Time {
	t := <-time.Tick(d)
	t.In(GetLocation())
	r := make(chan Time, 1)
	r <- t
	return r
}
func ParseDuration(s string) (Duration, error) {
	return time.ParseDuration(s)
}
func Sleep(d Duration) {
	time.Sleep(d)
}
func Parse(layout string, value string) (Time, error) {
	return time.ParseInLocation(layout, value, GetLocation())
}
func Unix(sec int64, nsec int64) Time {
	return time.Unix(sec, nsec).In(GetLocation())
}

func GetLocation() *Location {
	return location
}

func At(t Time) Time {
	return t.In(location)
}
