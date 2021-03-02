package util

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
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

func Float64ToStrWithDecimal(num float64, decimal int) string {
	d := float64(1)
	if decimal > 0 {
		d = math.Pow10(decimal)
	}
	return strconv.FormatFloat(math.Trunc(num*d)/d, 'f', -1, 64)
}

func ToJsonIgnoreErr(obj interface{}) string {
	resBytes, _ := json.Marshal(obj)
	return string(resBytes)
}

func ReadFile(filename string) (string, error) {
	if fileObj, err := os.Open(filename); err == nil {
		defer fileObj.Close()
		if contents, err := ioutil.ReadAll(fileObj); err == nil {
			return string(contents), nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}
