package util

import (
	"encoding/json"
	"strconv"
)

func StrToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func StrToInt64IgnoreErr(str string) int64 {
	i, _ := StrToInt64(str)
	return i
}

func Float64ToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func StrToFloat64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func ToJsonIgnoreErr(obj interface{}) string {
	resBytes, _ := json.Marshal(obj)
	return string(resBytes)
}
