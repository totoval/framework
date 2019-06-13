package bigfloat

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/totoval/framework/helpers/log"

	"github.com/totoval/framework/model/types/bigint"
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
type testDiv struct {
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
type testMarshalJSON struct {
	a      string
	output string
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
var testDivTable = []*testDiv{
	{"-123456789012345678901234567890.12345678901234567890123456789012345678901234567890", "1", "-123456789012345678901234567890.12345678901234567890123456789012345678901234567890"},
	{"123456789012345678901234567890.12345678901234567890123456789012345678901234567890", "1", "123456789012345678901234567890.12345678901234567890123456789012345678901234567890"},
	{"-12345678901234567890123456789012345678901234567890.123456789012345678901234567891", "1", "-12345678901234567890123456789012345678901234567890.123456789012345678901234567891"},
	{"-123456789012345678901234567890.123456789012345678901234567891", "1", "-123456789012345678901234567890.123456789012345678901234567891"},
	{"123.123", "1", "123.123"},
	{"-123.123", "1", "-123.123"},
	{"123123123.123456", "1", "123123123.123456"},
	{"-123123.123", "1", "-123123.123"},
	{"1", "1", "1"},
	{"1", "0.1", "10"},
	{"30000000000", "0.1", "300000000000"},
	{"30000000000.1235435", "0.1", "300000000001.235435"},
	{"123456", "123456", "1"},
	{"123456.1234", "123456.1234", "1"},
	{"123456.123400", "123456.1234", "1"},
	{"1000000000000000", "1000000000000000000", "0.001"},
	{"100000", "365", "273.97260274...."}, // infinitive
	{"123456789", "1000000000", "0.123456789"},
	{"1234567890123456789012345678901", "10000000000000000000000000000000", "0.1234567890123456789012345678901"},
	{"1234567890123456789012345678901", "1000000000000000000000000000", "1234.567890123456789012345678901"},
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
	{"-123456789012345678901234567890.0", 0, RoundUpAlways, "-123456789012345678901234567890"},
	{"-123456789012345678901234567890.0", 0, RoundUpAuto, "-123456789012345678901234567890"},
	{"-123456789012345678901234567890.0", 0, RoundDown, "-123456789012345678901234567890"},
	{"10.0", 0, RoundUpAlways, "10"},
	{"10.0", 0, RoundUpAuto, "10"},
	{"10.0", 0, RoundDown, "10"},
	{"-10.0", 0, RoundUpAlways, "-10"},
	{"-10.0", 0, RoundUpAuto, "-10"},
	{"-10.0", 0, RoundDown, "-10"},
	{"123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 4, RoundUpAuto, "123456789012345678901234567890.1235"},
	{"123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 4, RoundDown, "123456789012345678901234567890.1234"},
	{"123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 4, RoundUpAlways, "123456789012345678901234567890.1235"},
	{"-123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 4, RoundUpAuto, "-123456789012345678901234567890.1234"},
	{"-123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 4, RoundDown, "-123456789012345678901234567890.1235"},
	{"-123456789012345678901234567890.12345678901234567890123456789012345678901234567890", 4, RoundUpAlways, "-123456789012345678901234567890.1234"},
	// {"-12345678901234567890123456789012345678901234567890.123456789012345678901234567891", 10, "-12345678901234567890123456789012345678901234567889.123456789012345678901234567891"},
	// {"-123456789012345678901234567890.123456789012345678901234567891", 1, "-123456789012345678901234567889.123456789012345678901234567891"},
	// {"123.123", 10 "124.123"},
	// {"-123.123", 10, "-122.123"},
	// {"123123123.123456", 5, "123123124.123456"},
	// {"-123123.123", 2, "-123122.123"},
	// {"1", 1, "2"},
}
var testMarshalJSONTable = []*testMarshalJSON{
	{a: "100.123", output: `{"str1":"\"test\"","bf1":"100.123","bi1":"100"}`},
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

func TestBigFloat_Div(t *testing.T) {
	for _, te := range testDivTable {
		a := BigFloat{}
		a.CreateFromString(te.a, ToNearestEven)

		b := BigFloat{}
		// b.SetString(te.b)
		b.CreateFromString(te.b, ToNearestEven)

		c := BigFloat{}
		c.Div(a, b)
		//fmt.Println(newData.Acc(), newData.Text('f', 168))

		if strings.TrimRight(c.String(), "0") != strings.TrimRight(te.output, "0") {
			t.Errorf("div expected %s, got %s a:%v, b:%v", te.output, c.String(), te.a, te.b)
		}

	}
}

func TestBigFloat_Round(t *testing.T) {
	for _, te := range testRoundTable {
		a := BigFloat{}
		a.CreateFromString(te.a, ToNearestEven)

		result, err := a.Round(te.decimal, te.roundType)

		if err != nil {
			log.Info(err.Error())
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
	a, _, _ := big.ParseFloat("123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890.12345678901234567890123456789012345678901234567890123456789012345678901234567890", 10, AutoPrec, big.ToNearestEven)

	log.Info(string(a.Prec()))
	log.Info(a.Text('f', 1024))
	b := BigFloat{}
	log.Info(string(a.Prec()))
	b.Convert(a)
	fmt.Println(b.String())
	e, _ := b.Round(2, RoundUpAlways)
	fmt.Println(e.String())

	c := big.Int{}
	c.SetString("312343532485823940580923458934", 10)
	d := bigint.BigInt{}
	d.Convert(&c)
	fmt.Println(d.String())
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
func TestBigFloat_MarshalJSON(t *testing.T) {

	type testMarshalJSONStruct struct {
		Str string        `json:"str1,string"`
		Bf  BigFloat      `json:"bf1,string"`
		Bi  bigint.BigInt `json:"bi1,string"`
	}

	for _, te := range testMarshalJSONTable {
		var test testMarshalJSONStruct
		test.Str = "test"
		var bf BigFloat
		bf.CreateFromString(te.a, ToNearestEven)
		test.Bf = bf
		var bi bigint.BigInt
		bi.CreateFromString(te.a, 10)
		test.Bi = bi

		js, _ := json.Marshal(test)

		fmt.Println(js)

		if string(js) != te.output {
			t.Errorf("marshalJSON expected %s, got %s a:%v", te.output, string(js), te.a)
		}

	}
}
