package lcd

import (
	"github.com/irisnet/sdk-go/client/basic"
	"github.com/irisnet/sdk-go/util"
	"testing"
)

var (
	c LiteClient
)

func TestMain(m *testing.M) {
	baseClient := basic.NewClient("http://v2.irisnet-lcd.rainbow.one")
	c = NewClient(baseClient)
	m.Run()
}

func TestClient_QueryAccount(t *testing.T) {
	address := "iaa1cx79dle0acdystaqvva4m2n3ta95v2pp4aw9ve"
	if res, err := c.QueryAccount(address); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(res))
	}
}
