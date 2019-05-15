package null

import (
	"math/big"

	"github.com/totoval/framework/helpers/zone"
	"github.com/totoval/framework/model/types/bigfloat"
	"github.com/totoval/framework/model/types/bigint"
)

var nullString string = ""
var nullByte byte = 0
var nullTime zone.Time = zone.Unix(0, 0)

const nullNumber = 0
const nullStringNumber = "0"

var nullMapStringInterface map[string]interface{} = map[string]interface{}{nullString: nullNumber}
var nullFloat32 float32 = nullNumber
var nullFloat64 float64 = nullNumber
var nullInt int = nullNumber
var nullInt8 int8 = nullNumber
var nullInt16 int16 = nullNumber
var nullInt32 int32 = nullNumber
var nullInt64 int64 = nullNumber
var nullUint uint = nullNumber
var nullUint8 uint8 = nullNumber
var nullUint16 uint16 = nullNumber
var nullUint32 uint32 = nullNumber
var nullUint64 uint64 = nullNumber

func String() *string {
	return &nullString
}

func MapStringInterface() *map[string]interface{} {
	return &nullMapStringInterface
}

func Byte() *byte {
	return &nullByte
}

func Float64() *float64 {
	return &nullFloat64
}
func Float32() *float32 {
	return &nullFloat32
}

func Int() *int {
	return &nullInt
}
func Int64() *int64 {
	return &nullInt64
}
func Int32() *int32 {
	return &nullInt32
}
func Int16() *int16 {
	return &nullInt16
}
func Int8() *int8 {
	return &nullInt8
}

func Uint() *uint {
	return &nullUint
}
func Uint64() *uint64 {
	return &nullUint64
}
func Uint32() *uint32 {
	return &nullUint32
}
func Uint16() *uint16 {
	return &nullUint16
}
func Uint8() *uint8 {
	return &nullUint8
}

func BigInt() *bigint.BigInt {
	b := big.Int{}
	b.SetInt64(nullNumber)
	v := bigint.BigInt{}
	v.Convert(&b)
	return &v
}
func BigFloat() (*bigfloat.BigFloat, error) {
	v := bigfloat.BigFloat{}
	err := v.CreateFromString(nullStringNumber, bigfloat.ToNearestEven)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
func Time() *zone.Time {
	return &nullTime
}
