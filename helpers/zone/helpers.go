package zone

import (
    "time"

    "github.com/totoval/framework/config"
)

var location *time.Location

func init(){
    var err error
    location, err = time.LoadLocation(config.GetString("app.timezone"))
    if err != nil {
        panic(err)
    }
}

func At(t time.Time) time.Time {
    return t.In(location)
}