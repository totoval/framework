package bigfloat

import (
	"fmt"
	"log"
	"math/big"
	"strings"
	"testing"
)

type testAdd struct {
	a      string
	b      string
	output string
}
type testMul struct {
	a      string
	b      string
	output string
}
type testRound struct {
	a         string
	decimal   uint
	roundType RoundType
	output    string
}

var testAddTable = []*testAdd{
	{"-123456789012345678901234567890.12345678901234567890123456789012345678901234567890", "1", "-123456789012345678901234567889.12345678901234567890123456789012345678901234567890"},
	{"123456789012345678901234567890.12345678901234567890123456789012345678901234567890", "1", "123456789012345678901234567891.12345678901234567890123456789012345678901234567890"},
	{"-12345678901234567890123456789012345678901234567890.123456789012345678901234567891", "1", "-12345678901234567890123456789012345678901234567889.123456789012345678901234567891"},
	{"-123456789012345678901234567890.123456789012345678901234567891", "1", "-123456789012345678901234567889.123456789012345678901234567891"},
	{"123.123", "1", "124.123"},
	{"-123.123", "1", "-122.123"},
	{"123123123.123456", "1", "123123124.123456"},
	{"-123123.123", "1", "-123122.123"},
	{"1", "1", "2"},
	{"1", "0.1", "1.1"},
}
var testMulTable = []*testMul{
	{"-123456789012345678901234567890.12345678901234567890123456789012345678901234567890", "1", "-123456789012345678901234567890.12345678901234567890123456789012345678901234567890"},
	{"123456789012345678901234567890.12345678901234567890123456789012345678901234567890", "1", "123456789012345678901234567890.12345678901234567890123456789012345678901234567890"},
	{"-12345678901234567890123456789012345678901234567890.123456789012345678901234567891", "1", "-12345678901234567890123456789012345678901234567890.123456789012345678901234567891"},
	{"-123456789012345678901234567890.123456789012345678901234567891", "1", "-123456789012345678901234567890.123456789012345678901234567891"},
	{"123.123", "1", "123.123"},
	{"-123.123", "1", "-123.123"},
	{"123123123.123456", "1", "123123123.123456"},
	{"-123123.123", "1", "-123123.123"},
	{"1", "1", "1"},
	{"1", "0.1", "0.1"},
	{"30000000000", "0.1", "3000000000"},
	{"30000000000.1235435", "0.1", "3000000000.01235435"},
	{"123456", "123456", "15241383936"},
	{"123456.1234", "123456.1234", "15241414404.95602756"},
	{"123456.123400", "123456.1234", "15241414404.95602756"},
}
var testRoundTable = []*testRound{
	{"-123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 0, RoundDown, "-123456789012345678901234567891"},
	{"-123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 0, RoundUpAlways, "-123456789012345678901234567890"},
	{"-123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 0, RoundUpAuto, "-123456789012345678901234567890"},
	{"123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 0, RoundUpAuto, "123456789012345678901234567890"},
	{"123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 0, RoundDown, "123456789012345678901234567890"},
	{"123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 0, RoundUpAlways, "123456789012345678901234567891"},
	{"123456789012345678901234567890.0", 0, RoundUpAlways, "123456789012345678901234567890"},
	{"123456789012345678901234567890.0", 0, RoundUpAuto, "123456789012345678901234567890"},
	{"123456789012345678901234567890.0", 0, RoundDown, "123456789012345678901234567890"},
	{"10.0", 0, RoundUpAlways, "10"},
	{"10.0", 0, RoundUpAuto, "10"},
	{"10.0", 0, RoundDown, "10"},
	{"123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 4, RoundUpAuto, "123456789012345678901234567890.1235"},
	{"123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 4, RoundDown, "123456789012345678901234567890.1234"},
	{"123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 4, RoundUpAlways, "123456789012345678901234567890.1235"},
	// {"-12345678901234567890123456789012345678901234567890.123456789012345678901234567891", 10, "-12345678901234567890123456789012345678901234567889.123456789012345678901234567891"},
	// {"-123456789012345678901234567890.123456789012345678901234567891", 1, "-123456789012345678901234567889.123456789012345678901234567891"},
	// {"123.123", 10 "124.123"},
	// {"-123.123", 10, "-122.123"},
	// {"123123123.123456", 5, "123123124.123456"},
	// {"-123123.123", 2, "-123122.123"},
	// {"1", 1, "2"},
}

func TestBigFloat_Add(t *testing.T) {
	for _, te := range testAddTable {
		a := BigFloat{}
		a.CreateFromString(te.a, ToNearestEven)

		b := BigFloat{}
		// b.SetString(te.b)
		b.CreateFromString(te.b, ToNearestEven)

		c := BigFloat{}
		c.Add(a, b)
		// fmt.Println(newData.Acc(), newData.Text('f', 168))

		if strings.TrimRight(c.String(), "0") != strings.TrimRight(te.output, "0") {
			t.Errorf("add expected %s, got %s a:%v, b:%v", te.output, c.String(), te.a, te.b)
		}

		newC := BigFloat{}
		newC.CreateFromString(te.output, ToNearestEven)
		newA := BigFloat{}
		newA.Sub(newC, b)
		if strings.TrimRight(newA.String(), "0") != strings.TrimRight(a.String(), "0") {
			t.Errorf("sub expected %s, got %s a:%v, b:%v", a.String(), newA.String(), te.output, te.b)
		}

	}
}
func TestBigFloat_Mul(t *testing.T) {
	for _, te := range testMulTable {
		a := BigFloat{}
		a.CreateFromString(te.a, ToNearestEven)

		b := BigFloat{}
		// b.SetString(te.b)
		b.CreateFromString(te.b, ToNearestEven)

		c := BigFloat{}
		c.Mul(a, b)
		// fmt.Println(newData.Acc(), newData.Text('f', 168))

		if strings.TrimRight(c.String(), "0") != strings.TrimRight(te.output, "0") {
			t.Errorf("mul expected %s, got %s a:%v, b:%v", te.output, c.String(), te.a, te.b)
		}

		newC := BigFloat{}
		newC.CreateFromString(te.output, ToNearestEven)
		newA := BigFloat{}
		newA.Div(newC, b)
		if strings.TrimRight(newA.String(), "0") != strings.TrimRight(a.String(), "0") {
			t.Errorf("div expected %s, got %s a:%v, b:%v", a.String(), newA.String(), te.output, te.b)
		}

	}
}

func TestBigFloat_Round(t *testing.T) {
	for _, te := range testRoundTable {
		a := BigFloat{}
		a.CreateFromString(te.a, ToNearestEven)

		result, err := a.Round(te.decimal, te.roundType)

		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println(result.String())
		if result.String() != te.output {
			t.Errorf("expected %s, got %s a:%v, decimal:%v, roundType:%v", te.output, result.String(), te.a, te.decimal, te.roundType)
		}
	}
}
func TestBigFloat_Convert(t *testing.T) {
	//a := big.Float{}
	a, _, _ := big.ParseFloat("3.123456789012345678901234567890123456789012345678901234567890123456789012345678901", 10, 1024, big.ToNearestEven)

	log.Println(a.Prec())
	log.Println(a.Text('f', 1024))
	b := BigFloat{}
	b.CreateFromString(a.Text('f', 1024), ToNearestEven)
	fmt.Println(b.String())
}

func TestBigFloat_Floor(t *testing.T) {

	a := BigFloat{}
	a.CreateFromString("30000000000", ToNearestEven)
	b := BigFloat{}
	b.CreateFromString("0.1", ToNearestEven)
	c := BigFloat{}
	c.Mul(a, b)
	fmt.Println(c.String())
	d, err := c.Floor()

	if d.String() != "3000000000" {
		t.Errorf("expected %s, got %s a:%v b:%v", "3000000000", d.String(), a.String(), b.String())
	}

	fmt.Println(d.String() == "10", err)
}
