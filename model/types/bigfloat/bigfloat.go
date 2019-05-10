package bigfloat

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/totoval/framework/model/types/bigint"
)

// These constants define supported rounding modes.
const (
	ToNearestEven big.RoundingMode = iota // == IEEE 754-2008 roundTiesToEven
	ToNearestAway                         // == IEEE 754-2008 roundTiesToAway
	ToZero                                // == IEEE 754-2008 roundTowardZero
	AwayFromZero                          // no IEEE 754-2008 equivalent
	ToNegativeInf                         // == IEEE 754-2008 roundTowardNegative
	ToPositiveInf                         // == IEEE 754-2008 roundTowardPositive
)

type _bf = big.Float
type BigFloat struct {
	_bf
	normalCount  uint
	decimalCount uint
}

func (bf *BigFloat) BF() _bf {
	return bf._bf
}

var ZERO BigFloat

func init() {
	if err := ZERO.CreateFromString("0", ToNearestEven); err != nil {
		panic(err)
	}
}

func (bf BigFloat) Value() (driver.Value, error) {
	//debug.Dump(bf._bf.Prec(), bf.Text('f', 100), bf.String())
	return []byte(bf.String()), nil
}
func (bf *BigFloat) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return bf.scanBytes(src)
	case string:
		return bf.scanBytes([]byte(src))
	case nil:
		bf = nil
		return nil
	default:
		return fmt.Errorf("pq: cannot convert %T to BigFloat", src)
	}

}

func (bf *BigFloat) scanBytes(src []byte) error {
	return bf.CreateFromString(string(src), ToNearestEven)
}
func (bf *BigFloat) String() string {
	//debug.Dump(bf._bf.Prec(), bf._bf.MinPrec())
	//if bf.decimalCount == 0 {
	//	return bf.Text('f', 62)
	//}
	return bf.Text('f', int(bf.decimalCount)/2)
}

func (bf *BigFloat) SetInt(i *bigint.BigInt, mode big.RoundingMode) error {
	return bf.CreateFromString(i.String(), mode)
}

func (bf *BigFloat) SetDecimal(d uint) {
	bf.decimalCount = d * 2
}

func (bf *BigFloat) CreateFromString(s string, mode big.RoundingMode) error {
	//parse number string
	parts := strings.Split(s, ".")
	if len(parts) == 1 {
		// There is no decimal point, we can just parse the original string as
		// an int
		bf.normalCount = uint(len(parts[0])) * 2
		bf.SetDecimal(0)
	} else if len(parts) == 2 {
		// strip the insignificant digits for more accurate comparisons.
		decimalPart := strings.TrimRight(parts[1], "0")
		bf.normalCount = uint(len(parts[0])) * 2
		bf.SetDecimal(uint(len(decimalPart)))
	} else {
		return errors.New("can't convert " + s + " to decimal")
	}

	// string to BigFloat
	_bf, _, err := big.ParseFloat(s, 10, bf.normalCount*2+bf.decimalCount*2+8, mode)
	bf._bf = *_bf
	//bf.SetPrec(prec).SetMode(mode)
	//_, err := fmt.Sscan(s, &bf._bf)
	return err
}

//@todo xml protobuf ...
func (bf *BigFloat) MarshalJSON() ([]byte, error) {
	return []byte(bf.String()), nil
}

func (bf *BigFloat) useBiggerDecimal(a BigFloat, b BigFloat) {
	if a.decimalCount > b.decimalCount {
		bf.decimalCount = a.decimalCount
	} else {
		bf.decimalCount = b.decimalCount
	}
}

func (bf *BigFloat) Add(a BigFloat, b BigFloat) {
	bf.useBiggerDecimal(a, b)
	bf._bf.Add(&a._bf, &b._bf)
}
func (bf *BigFloat) Sub(a BigFloat, b BigFloat) {
	bf.useBiggerDecimal(a, b)
	bf._bf.Sub(&a._bf, &b._bf)
}
func (bf *BigFloat) Mul(a BigFloat, b BigFloat) {
	bf.useBiggerDecimal(a, b)
	bf._bf.Mul(&a._bf, &b._bf)
}
func (bf *BigFloat) Div(a BigFloat, b BigFloat) {
	bf.useBiggerDecimal(a, b)
	bf._bf.Quo(&a._bf, &b._bf)
}
func (bf *BigFloat) Abs(a BigFloat) {
	bf._bf.Abs(&a._bf)
}
func (bf *BigFloat) Cmp(a BigFloat) int {
	return bf._bf.Cmp(&a._bf)
}

//
//func main(){
//	a := BigFloat{}
//	a.SetString("10", 10)
//	b := BigFloat{}
//	b.SetString("11", 10)
//	c := BigFloat{}
//	c.Add(&a.BF, &b.BF)
//}
