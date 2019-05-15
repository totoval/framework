package pb

import (
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/totoval/framework/helpers/zone"
)

func DurationConvert(d *duration.Duration) *zone.Duration {
	result := zone.Duration(d.GetSeconds())*zone.Second + zone.Duration(d.GetNanos())*zone.Nanosecond
	return &result
}
func TimestampConvert(t *timestamp.Timestamp) *zone.Time {
	result := zone.Unix(t.GetSeconds(), int64(t.GetNanos()))
	return &result
}
