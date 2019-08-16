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

type BigFloat struct {
	_bf          big.Float
	normalCount  uint
	decimalCount uint
}

const AutoPrec = 512 // 256 -> decimal 32   512 -> decimal 78

func Zero() *BigFloat {
	zero := &BigFloat{}
	_ = zero.CreateFromString("0", ToNearestEven)
	return zero
}

func (bf *BigFloat) Convert(f *big.Float) error {
	// int(f.Prec()) uint to int may cause precision loss
	prec := f.Prec()
	if prec > big.MaxExp {
		return errors.New("precision is too large, may cause precision loss")
	}
	return bf.CreateFromString(f.Text('f', int(prec)), ToNearestEven)
}
func (bf *BigFloat) Float() *big.Float {
	return &bf._bf
}

func (bf *BigFloat) BF() big.Float {
	return bf._bf
}

func (bf BigFloat) Value() (driver.Value, error) {
	// debug.Dump(bf._bf.Prec(), bf.Text('f', 100), bf.String())
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
	// result := bf.Text('f', int(bf.Prec()))
	//
	// switch bf.Acc() {
	// case big.Above:
	// 	for i := bf.Prec(); i > 0; i-- {
	// 		result = bf.Text('f', int(i))
	// 		if bf.Acc() == big.Exact {
	// 			break
	// 		}
	// 	}
	// 	break
	// case big.Below:
	// 	for i := uint(0); i <= bf.Prec(); i++ {
	// 		result = bf.Text('f', int(i))
	// 		if bf.Acc() == big.Exact {
	// 			break
	// 		}
	// 	}
	// 	break
	// case big.Exact:
	// 	break
	// }
	//
	// trimResult := strings.TrimRight(result, "0")
	//
	// if trimResult[len(trimResult)-1:] == "." {
	// 	trimResult = trimResult[:len(trimResult)-1]
	// }
	//
	// return trimResult

	result := bf._bf.Text('f', int(bf.decimalCount)/2)
	trimResult := result
	if strings.Contains(result, ".") {
		trimResult = strings.TrimRight(result, "0")
		if trimResult[len(trimResult)-1:] == "." {
			trimResult = trimResult[:len(trimResult)-1]
		}
	}

	return trimResult
}

func (bf *BigFloat) SetInt(i *bigint.BigInt, mode big.RoundingMode) error {
	return bf.CreateFromString(i.String(), mode)
}

func (bf *BigFloat) setDecimal(d uint) { // @todo 0 is infinity
	bf.decimalCount = d * 2
}

func (bf *BigFloat) Copy(newBf *BigFloat) error {
	return newBf.CreateFromString(bf.String(), bf._bf.Mode())
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

	// tmp := ""
	// if lastDecimal == 0{
	// 	tmp = newDecimalPartPlusStr
	// }else{
	// 	tmp =
	// }
	// newDecimalPartPlusStr[:len(newDecimalPartPlusStr)-1]

	// var newDecimalPartPlus BigFloat
	// err := newDecimalPartPlus.CreateFromString(newDecimalPartPlusStr, ToNearestEven)
	// if err != nil {
	// 	return nil, err
	// }

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

func (bf *BigFloat) Round(decimal uint, roundType RoundType) (*BigFloat, error) {
	var bfCopy BigFloat
	if err := bf.Copy(&bfCopy); err != nil {
		return nil, err
	}
	parts := strings.Split(bfCopy.String(), ".")
	normalPart := ""
	decimalPart := ""
	if len(parts) == 1 {
		normalPart = parts[0]
		decimalPart = ""
		bfCopy.setDecimal(0)
	} else if len(parts) == 2 {
		normalPart = parts[0]
		decimalPart = parts[1]
	} else {
		return nil, errors.New("cannot parse " + bfCopy.String())
	}

	// if provide decimal is greater than the real decimal, then there isn't any precision problem, so directly return
	if int(decimal) >= len(decimalPart) {
		return &bfCopy, nil
	}

	result := &BigFloat{}
	var err error
	// check is greater than 0
	if bfCopy.Cmp(*Zero()) < 0 {
		//return nil, errors.New("currently not support for number smaller than 0")
		//@todo small than 0
		result, err = smallerThanZero(decimalPart, normalPart, decimal, roundType)
	} else {
		result, err = greaterOrEqualThanZero(decimalPart, normalPart, decimal, roundType)
	}
	if err != nil {
		return nil, err
	}

	// result.setDecimal(decimal)
	if err := result.CreateFromString(result.String(), ToNearestEven); err != nil {
		return nil, err
	}
	return result, nil
}

func greaterOrEqualThanZero(decimalPart string, normalPart string, decimal uint, roundType RoundType) (*BigFloat, error) {
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

	return result, nil
}
func smallerThanZero(decimalPart string, normalPart string, decimal uint, roundType RoundType) (*BigFloat, error) {
	//debug.Dump(normalPart, decimalPart, decimal, roundType)
	// -123.12345
	// decimal 0
	newDecimalPart := decimalPart[:decimal]            // ""
	lastDecimalStr := decimalPart[decimal : decimal+1] // 1
	//debug.Dump(lastDecimalStr, newDecimalPart)
	lastDecimal, err := strconv.ParseUint(lastDecimalStr, 10, 32) // 1
	if err != nil {
		return nil, err
	}
	newDecimalPartPlus := newDecimalPart + lastDecimalStr // 1

	// create roundDownPlus with RoundDown decimal + 1              decimal = 2         1000.1234 => 1000.123
	roundUpPlusStr := normalPart + "." + newDecimalPartPlus // -123.1
	var roundUpPlus BigFloat
	if err := roundUpPlus.CreateFromString(roundUpPlusStr, ToNearestEven); err != nil {
		return nil, err
	}
	//debug.Dump(roundDownPlusStr)

	// create roundUp with RoundUp                                 decimal = 2         1000.123 => 1000.12
	roundUpStr := normalPart + "." + newDecimalPart // -123.
	var roundUp BigFloat
	if err := roundUp.CreateFromString(roundUpStr, ToNearestEven); err != nil {
		return nil, err
	}
	//debug.Dump(roundUp)
	// create carry
	carry, err := createCarry(uint(lastDecimal), newDecimalPartPlus)
	// debug.Dump(lastDecimal, newDecimalPartPlus, carry)
	if err != nil {
		return nil, err
	}

	result := &BigFloat{}
	switch roundType {
	case RoundUpAlways:
		result = &roundUp
		break
	case RoundUpAuto:
		if lastDecimal <= 5 {
			result = &roundUp
		} else {
			result.Sub(roundUpPlus, *carry)
		}
		break
	case RoundDown:
		if lastDecimal > 0 {
			result.Sub(roundUpPlus, *carry)
		} else {
			result = &roundUp
		}
		break

	default:
		return nil, errors.New("unknown roundType")
	}

	return result, nil
}

func (bf *BigFloat) Ceil() (*BigFloat, error) {
	return bf.Round(0, RoundUpAlways)
}
func (bf *BigFloat) Floor() (*BigFloat, error) {
	return bf.Round(0, RoundDown)
}

func (bf *BigFloat) CreateFromString(s string, mode big.RoundingMode) error {
	// parse number string
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
	// _bf, _, err := big.ParseFloat(s, 10, bf.normalCount*2+bf.decimalCount*2+8, mode)
	_bf, _, err := big.ParseFloat(s, 10, AutoPrec, mode)
	// _bf, _, err := big.ParseFloat(s, 10, 2, mode)
	if err != nil {
		return err
	}
	bf._bf = *_bf
	return nil

	// tmp := &big.Float{}
	// // _, _, err := tmp.Parse(s, 10)
	// // tmp, _, err := big.ParseFloat(s, 10, bf.normalCount*2+bf.decimalCount*2+8, mode)
	// tmp, _, err := big.ParseFloat(s, 10, 168, mode)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(tmp.Acc())
	// bf._bf = *tmp

	// bf.SetPrec(prec).SetMode(mode)
	// _, err := fmt.Sscan(s, &bf._bf)
	// return err
}

// @todo xml protobuf ...
func (bf BigFloat) MarshalJSON() ([]byte, error) {
	// fix https://github.com/golang/go/issues/20651
	return []byte(`"` + bf.String() + `"`), nil
}

func (bf *BigFloat) UnmarshalJSON(src []byte) error {
	return bf.scanBytes(src)
}

func (bf *BigFloat) UnmarshalBinary(data []byte) error {
	return bf.scanBytes(data)
}

func (bf BigFloat) MarshalBinary() (data []byte, err error) {
	return []byte(bf.String()), nil
}
func (bf *BigFloat) useBiggerDecimal(a BigFloat, b BigFloat) {
	if a.decimalCount > b.decimalCount {
		bf.decimalCount = a.decimalCount
	} else {
		bf.decimalCount = b.decimalCount
	}
	if a.normalCount > b.normalCount {
		bf.normalCount = a.normalCount
	} else {
		bf.normalCount = b.normalCount
	}
}

func (bf *BigFloat) mergeDecimal(a BigFloat, b BigFloat) {
	bf.decimalCount = a.decimalCount + b.decimalCount
}
func (bf *BigFloat) mergeDecimalDiv(a BigFloat, b BigFloat, isInf ...bool) {
	decimalA := a.decimalCount
	decimalB := b.decimalCount

	if len(isInf) > 0 {
		if isInf[0] {
			bf.decimalCount = AutoPrec / 2
			return
		}
	}

	if decimalA == 0 && decimalB == 0 {
		// may be infinitive
		bf.decimalCount = AutoPrec / 2
		return
	}

	if decimalA == 0 {
		decimalA = 1
	}
	if decimalB == 0 {
		decimalB = 1
	}

	bf.decimalCount = decimalA * decimalB
	return
}

//@todo calc pointer param
func (bf *BigFloat) Add(a BigFloat, b BigFloat) {
	bf.useBiggerDecimal(a, b)
	bf._bf.Add(&a._bf, &b._bf)
}
func (bf *BigFloat) Sub(a BigFloat, b BigFloat) {
	bf.useBiggerDecimal(a, b)
	bf._bf.Sub(&a._bf, &b._bf)
}
func (bf *BigFloat) Mul(a BigFloat, b BigFloat) {
	bf.mergeDecimal(a, b)
	bf._bf.Mul(&a._bf, &b._bf)
}
func (bf *BigFloat) Div(a BigFloat, b BigFloat, isInf ...bool) {
	bf.mergeDecimalDiv(a, b, isInf...)
	bf._bf.Quo(&a._bf, &b._bf)
}
func (bf *BigFloat) Abs(a BigFloat) {
	bf._bf.Abs(&a._bf)
}
func (bf *BigFloat) Cmp(a BigFloat) int {
	return bf._bf.Cmp(&a._bf)
}

//
// func main(){
// 	a := BigFloat{}
// 	a.SetString("10", 10)
// 	b := BigFloat{}
// 	b.SetString("11", 10)
// 	c := BigFloat{}
// 	c.Add(&a.BF, &b.BF)
// }
