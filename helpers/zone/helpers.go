package zone

import (
	"time"

	"github.com/totoval/framework/config"
)

var location *time.Location

type Time = time.Time
type Duration = time.Duration
type Location = time.Location

func init() {
	var err error
	location, err = time.LoadLocation(config.GetString("app.timezone"))
	if err != nil {
		panic(err)
	}
}

func Now() Time {
	return time.Now().In(GetLocation())
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
