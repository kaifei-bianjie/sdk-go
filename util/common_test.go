package util

import (
	"math"
	"testing"
)

func TestFloat64ToStr(t *testing.T) {
	v := float64(23214.234543)
	v1 := v * math.Pow10(18)
	res := Float64ToStr(v1)
	t.Log(res)
}

func TestStrToInt64(t *testing.T) {

}
