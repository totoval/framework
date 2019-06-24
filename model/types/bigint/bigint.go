package bigint

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math/big"
)

type BigInt struct {
	_bi big.Int
}

func (bi *BigInt) BI() big.Int {
	return bi._bi
}

func Zero() *BigInt {
	zero := &BigInt{}
	_ = zero.CreateFromString("0", 10)
	return zero
}

func (bi *BigInt) Convert(i *big.Int) error {
	return bi.CreateFromString(i.Text(10), 10)
}

func (bi *BigInt) Int() *big.Int {
	return &bi._bi
}

func (bi BigInt) String() string {
	return bi._bi.String()
}

func (bi BigInt) Value() (driver.Value, error) {
	return []byte(bi._bi.String()), nil
}
func (bi *BigInt) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return bi.scanBytes(src)
	case string:
		return bi.scanBytes([]byte(src))
	case nil:
		bi = nil
		return nil
	default:
		return fmt.Errorf("pq: cannot convert %T to BigInt", src)
	}
}

func (bi *BigInt) scanBytes(src []byte) error {
	return bi.CreateFromString(string(src), 10)
}

func (bi *BigInt) CreateFromString(s string, base int) error {
	_, ok := bi._bi.SetString(s, base)
	if !ok {
		return errors.New("create bigint from string failed: " + s)
	}
	return nil
}

// @todo xml protobuf ...
func (bi BigInt) MarshalJSON() ([]byte, error) {
	// fix https://github.com/golang/go/issues/20651
	return []byte(`"` + bi._bi.String() + `"`), nil
}
func (bi *BigInt) UnmarshalJSON(src []byte) error {
	return bi.scanBytes(src)
}

func (bi *BigInt) UnmarshalBinary(data []byte) error {
	return bi.scanBytes(data)
}

func (bi BigInt) MarshalBinary() (data []byte, err error) {
	return []byte(bi.String()), nil
}

func (bi *BigInt) SetUint64(i uint64) *BigInt {
	bi._bi.SetUint64(i)
	return bi
}
func (bi *BigInt) SetInt64(i int64) *BigInt {
	bi._bi.SetInt64(i)
	return bi
}

//@todo calc pointer param
func (bi *BigInt) Add(a BigInt, b BigInt) {
	bi._bi.Add(&a._bi, &b._bi)
}
func (bi *BigInt) Sub(a BigInt, b BigInt) {
	bi._bi.Sub(&a._bi, &b._bi)
}
func (bi *BigInt) Mul(a BigInt, b BigInt) {
	bi._bi.Mul(&a._bi, &b._bi)
}
func (bi *BigInt) Div(a BigInt, b BigInt) {
	bi._bi.Quo(&a._bi, &b._bi)
}
func (bi *BigInt) Pow(a BigInt, b BigInt) error {
	if b.Cmp(*Zero()) < 0 {
		return errors.New("b cannot be smaller than 0")
	}
	bi._bi.Exp(&a._bi, &b._bi, nil)
	return nil
}
func (bi *BigInt) Abs(a BigInt) {
	bi._bi.Abs(&a._bi)
}
func (bi *BigInt) Cmp(a BigInt) int {
	return bi._bi.Cmp(&a._bi)
}

//
// func main(){
// 	a := BigInt{}
// 	a.SetString("10", 10)
// 	b := BigInt{}
// 	b.SetString("11", 10)
// 	c := BigInt{}
// 	c.Add(&a.BI, &b.BI)
// }
