package ad

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/irisnet/sdk-go/util"
	"strings"
	"testing"
)

func TestReadFromExcel(t *testing.T) {
	xlsx, err := excelize.OpenFile("./Sim-Upgrade-Validators-formatted.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	multipleAddrMap := make(map[string]int64)
	addrAmtMap := make(map[string]float64)
	sheet := "Sheet1"
	rows := xlsx.GetRows(sheet)
	for i := 0; i < len(rows); i++ {
		// set cell axis
		addrAxis := fmt.Sprintf("A%d", i+2)
		amtAxis := fmt.Sprintf("B%d", i+2)

		// get cell value
		addr := strings.ToLower(xlsx.GetCellValue(sheet, addrAxis))
		addr = strings.Replace(addr, " ", "", -1)
		addr = strings.Replace(addr, "\n", "", -1)
		addr = strings.Replace(addr, "\t", "", -1)
		amtStr := xlsx.GetCellValue(sheet, amtAxis)

		if addr == "" || amtStr == "" ||
			addr == "iaa1fd8x2ww89ztq7qdrhw8x3z5aj25svxy3588n6u" ||
			addr == "iaa180z3qagykwpr7v6htawvh7z3n5t7zw6w0zjvc2" {
			continue
		}

		// convert str to float64
		var amt float64
		if v, err := util.StrToFloat64(amtStr); err != nil {
			t.Fatalf("convert amt to float fail, row: %d, err: %s\n", i, err)
		} else {
			amt = v
		}

		if v, ok := addrAmtMap[addr]; ok {
			addrAmtMap[addr] = v + amt

			if v1, ok1 := multipleAddrMap[addr]; ok1 {
				multipleAddrMap[addr] = v1 + 1
			} else {
				multipleAddrMap[addr] = 1
			}
		} else {
			addrAmtMap[addr] = amt
		}
	}

	var count int64
	for _, v := range multipleAddrMap {
		count += v
	}

	t.Logf("origin sheet length: %d, handled result length:%d\n", len(rows), len(addrAmtMap))
	t.Logf("multipleAddrs length: %d\n", count)
	t.Log(util.ToJsonIgnoreErr(addrAmtMap))
	t.Log(util.ToJsonIgnoreErr(multipleAddrMap))
}
