package bigfloat

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math/big"
	"strconv"
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

func (bf *BigFloat) setDecimal(d uint) {
	bf.decimalCount = d * 2
}

func (bf *BigFloat) Copy(newBf *BigFloat) error {
	return newBf.CreateFromString(bf.String(), bf.Mode())
}

type RoundType byte

const (
	RoundUpAlways RoundType = iota
	RoundDown
	RoundUpAuto
)

func createCarry(lastDecimal uint, newDecimalPartPlusStr string) (*BigFloat, error) {
	decimal := len(newDecimalPartPlusStr)

	carryLastDecimal := uint(0)
	if lastDecimal > 0 {
		carryLastDecimal = 10 - lastDecimal
	} else {
		carryLastDecimal = 0
	}

	//tmp := ""
	//if lastDecimal == 0{
	//	tmp = newDecimalPartPlusStr
	//}else{
	//	tmp =
	//}
	//newDecimalPartPlusStr[:len(newDecimalPartPlusStr)-1]

	//var newDecimalPartPlus BigFloat
	//err := newDecimalPartPlus.CreateFromString(newDecimalPartPlusStr, ToNearestEven)
	//if err != nil {
	//	return nil, err
	//}

	carryStr := "0."
	for i := 0; i < decimal; i++ {
		if i == decimal-1 {
			carryStr += fmt.Sprintf("%d", carryLastDecimal)
		} else {
			carryStr += "0"
		}
	}
	var carry BigFloat
	if err := carry.CreateFromString(carryStr, ToNearestEven); err != nil {
		return nil, err
	}
	return &carry, nil
}

//func (bf *BigFloat) roundDown(decimal uint) (*BigFloat, error) {
//	var tmp BigFloat
//	if err := bf.Copy(&tmp); err != nil {
//		return nil, err
//	}
//	parts := strings.Split(tmp.String(), ".")
//	normalPart := parts[0]
//	decimalPart := parts[1]
//
//	// if provide decimal is greater than the real decimal, then there isn't any precision problem, so directly return
//	if int(decimal) > len(decimalPart) {
//		return bf, nil
//	}
//
//	newDecimalPart := decimalPart[:decimal]
//	lastDecimal, err := strconv.ParseUint(decimalPart[decimal:decimal+1], 10, 32)
//	if err != nil {
//		return nil, err
//	}
//
//	// create roundDown with RoundDown
//	roundDownStr := normalPart + "." + newDecimalPart
//	var roundDown BigFloat
//	if err := roundDown.CreateFromString(roundDownStr, ToNearestEven); err != nil {
//		return nil, err
//	}
//}
func (bf *BigFloat) Round(decimal uint, roundType RoundType) (*BigFloat, error) {
	var tmp BigFloat
	if err := bf.Copy(&tmp); err != nil {
		return nil, err
	}
	parts := strings.Split(tmp.String(), ".")
	normalPart := parts[0]
	decimalPart := parts[1]

	// check is greater than 0
	if tmp.Cmp(ZERO) < 0 {
		return nil, errors.New("currently not support for number smaller than 0")
	}

	// if provide decimal is greater than the real decimal, then there isn't any precision problem, so directly return
	if int(decimal) > len(decimalPart) {
		return bf, nil
	}

	newDecimalPart := decimalPart[:decimal]
	lastDecimalStr := decimalPart[decimal : decimal+1]
	lastDecimal, err := strconv.ParseUint(lastDecimalStr, 10, 32)
	if err != nil {
		return nil, err
	}
	newDecimalPartPlus := newDecimalPart + lastDecimalStr

	// create roundDownPlus with RoundDown decimal + 1              decimal = 2         1000.1234 => 1000.123
	roundDownPlusStr := normalPart + "." + newDecimalPartPlus
	var roundDownPlus BigFloat
	if err := roundDownPlus.CreateFromString(roundDownPlusStr, ToNearestEven); err != nil {
		return nil, err
	}

	// create roundDown with RoundDown                                 decimal = 2         1000.123 => 1000.12
	roundDownStr := normalPart + "." + newDecimalPart
	var roundDown BigFloat
	if err := roundDown.CreateFromString(roundDownStr, ToNearestEven); err != nil {
		return nil, err
	}

	// create carry
	carry, err := createCarry(uint(lastDecimal), newDecimalPartPlus)
	if err != nil {
		return nil, err
	}

	result := &BigFloat{}
	switch roundType {
	case RoundUpAlways:
		if lastDecimal > 0 {
			result.Add(roundDownPlus, *carry)
		} else {
			result = &roundDown
		}
		break
	case RoundUpAuto:
		if lastDecimal >= 5 {
			result.Add(roundDownPlus, *carry)
		} else {
			result = &roundDown
		}
		break
	case RoundDown:
		result = &roundDown
		break

	default:
		return nil, errors.New("unknown roundType")
	}

	result.setDecimal(decimal)
	return result, nil
}

func (bf *BigFloat) Ceil() (*BigFloat, error) {
	return bf.Round(0, RoundUpAlways)
}
func (bf *BigFloat) Floor() (*BigFloat, error) {
	return bf.Round(0, RoundDown)
}

func (bf *BigFloat) CreateFromString(s string, mode big.RoundingMode) error {
	//parse number string
	parts := strings.Split(s, ".")
	if len(parts) == 1 {
		// There is no decimal point, we can just parse the original string as
		// an int
		bf.normalCount = uint(len(parts[0])) * 2
		bf.setDecimal(0)
	} else if len(parts) == 2 {
		// strip the insignificant digits for more accurate comparisons.
		decimalPart := strings.TrimRight(parts[1], "0")
		bf.normalCount = uint(len(parts[0])) * 2
		bf.setDecimal(uint(len(decimalPart)))
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
