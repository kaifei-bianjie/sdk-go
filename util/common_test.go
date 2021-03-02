package util

import (
	"fmt"
	"math"
	"testing"
)

func TestFloat64ToStr(t *testing.T) {
	v := float64(0.0000677806639676961)
	v1 := v * math.Pow10(18)
	res := Float64ToStr(v1)
	t.Log(res)
}

func TestFloat64ToStrWithDecimal(t *testing.T) {
	v := float64(0.0000677806639676961)
	res := Float64ToStrWithDecimal(v, 10)
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
	amtStr := "0.0000677806639676961"
	amtF, err := StrToFloat64(amtStr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("%.8f", amtF))
}
