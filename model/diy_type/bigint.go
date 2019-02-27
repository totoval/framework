package diy_type

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/totoval/framework/helpers"
	"math/big"
)

type BI = big.Int
type BigInt struct {
	BI
}

func (bi BigInt) Value() (driver.Value, error) {
	return []byte(bi.String()), nil
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
		helpers.DD(src)
	}

	return fmt.Errorf("pq: cannot convert %T to BoolArray", src)
}

func (bi *BigInt) scanBytes(src []byte) error {
	return bi.CreateFromString(string(src), 10)
}

func (bi *BigInt) CreateFromString(s string, base int) error {
	_, ok := bi.SetString(s, base)
	if !ok {
		return errors.New("create bigint from string failed: "+ s)
	}
	return nil
}

func (bi *BigInt) Add(a BigInt, b BigInt) {
	bi.BI.Add(&a.BI, &b.BI)
}
func (bi *BigInt) Sub(a BigInt, b BigInt) {
	bi.BI.Sub(&a.BI, &b.BI)
}
func (bi *BigInt) Mul(a BigInt, b BigInt) {
	bi.BI.Mul(&a.BI, &b.BI)
}
func (bi *BigInt) Div(a BigInt, b BigInt) {
	bi.BI.Quo(&a.BI, &b.BI)
}
func (bi *BigInt) Abs(a BigInt) {
	bi.BI.Abs(&a.BI)
}
func (bi *BigInt) Cmp (a BigInt) int {
	return bi.BI.Cmp(&a.BI)
}

//
//func main(){
//	a := BigInt{}
//	a.SetString("10", 10)
//	b := BigInt{}
//	b.SetString("11", 10)
//	c := BigInt{}
//	c.Add(&a.BI, &b.BI)
//}