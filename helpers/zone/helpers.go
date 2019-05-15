package zone

import (
	"time"

	"github.com/totoval/framework/config"
)

var location *time.Location

func init() {
	var err error
	location, err = time.LoadLocation(config.GetString("app.timezone"))
	if err != nil {
		panic(err)
	}
}

func Now() time.Time {
	return time.Now().In(Location())
}
func Parse(layout string, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, Location())
}
func Unix(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec).In(Location())
}

func Location() *time.Location {
	return location
}

func At(t time.Time) time.Time {
	return t.In(location)
}
