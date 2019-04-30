package pb

import (
	"time"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func DurationConvert(d *duration.Duration) *time.Duration {
	result := time.Duration(d.GetSeconds())*time.Second + time.Duration(d.GetNanos())*time.Nanosecond
	return &result
}
func TimestampConvert(t *timestamp.Timestamp) *time.Time {
	result := time.Unix(t.GetSeconds(), int64(t.GetNanos()))
	return &result
}
