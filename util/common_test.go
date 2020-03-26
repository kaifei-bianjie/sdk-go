package util

import (
	"math"
	"strings"
	"testing"
)

func TestFloat64ToStr(t *testing.T) {
	v := float64(23214.234543)
	v1 := v * math.Pow10(18)
	res := Float64ToStr(v1)
	t.Log(res)
}

func TestStrToFloat64(t *testing.T) {
	str := "3.1415926535"
	if v, err := StrToFloat64(str); err != nil {
		t.Fatal(err)
	} else {
		t.Log(v)
	}
}

func TestStr(t *testing.T) {
	t.Log(strings.Replace(" bnb18krp3ktjgum0tcdd23d6wn587qlsg54llk6kjs", " ", "", -1))
}
