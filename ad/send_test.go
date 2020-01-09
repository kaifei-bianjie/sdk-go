package ad

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/sdk-go/client/basic"
	"github.com/irisnet/sdk-go/client/lcd"
	"github.com/irisnet/sdk-go/client/rpc"
	"github.com/irisnet/sdk-go/client/tx"
	ctypes "github.com/irisnet/sdk-go/client/types"
	"github.com/irisnet/sdk-go/keys"
	"github.com/irisnet/sdk-go/types"
	"github.com/irisnet/sdk-go/util"
	"math"
	"testing"
	"time"
)

type (
	MnemonicKMParams struct {
		Menmonic string
		Password string
		FullPath string
	}

	KeyStoreKMParams struct {
		FilePath string
		Password string
	}
)

var (
	liteClient lcd.LiteClient
	rpcClient  rpc.RPCClient
	txClient   tx.TxClient
	km         keys.KeyManager
)

func initClient() {
	//network := types.Mainnet
	//lcdUrl := "http://v2.irisnet-lcd.rainbow.one"
	//rpcUrl := "tcp://seed-1.mainnet.irisnet.org:26657"
	//initKMType := "seed"

	network := types.Testnet
	lcdUrl := "http://irisnet-lcd.dev.bianjie.ai"
	rpcUrl := "tcp://192.168.150.32:26657"
	initKMType := "seed"

	switch initKMType {
	case "seed":
		p := MnemonicKMParams{
			Menmonic: "",
			Password: "",
			FullPath: "44'/118'/1'/0/0",
		}
		if v, err := keys.NewMnemonicKeyManager(p.Menmonic, p.Password, p.FullPath); err != nil {
			fmt.Println("init km fail")
			panic(err)
		} else {
			km = v
		}
		break
	case "ks":
		p := KeyStoreKMParams{
			FilePath: "",
			Password: "",
		}
		if v, err := keys.NewKeyStoreKeyManager(p.FilePath, p.Password); err != nil {
			fmt.Println("init km fail")
			panic(err)
		} else {
			km = v
		}
		break
	default:
		panic("should init km first")
	}

	switch network {
	case types.Mainnet:
		sdk.SetNetworkType("mainnet")
	case types.Testnet:
		sdk.SetNetworkType("testnet")
	}

	fmt.Printf("address is: %s\n", km.GetAddr().String())

	basicClient := basic.NewClient(lcdUrl)
	liteClient = lcd.NewClient(basicClient)
	rpcClient = rpc.NewClient(rpcUrl)
	var chainId string
	if status, err := rpcClient.GetStatus(); err != nil {
		fmt.Println("init rpc client err")
		panic(err)
	} else {
		chainId = status.NodeInfo.Network
	}

	if c, err := tx.NewClient(chainId, network, km, liteClient, rpcClient); err != nil {
		fmt.Println("init tx client err")
		panic(err)
	} else {
		txClient = c
	}
}

func TestInitClient(m *testing.T) {
	initClient()
}

func TestGetAccountInfo(t *testing.T) {
	initClient()
	address := "faa1d4y0pw87s05f93mtt88fad9sag4fna5az08kfu"
	if v, err := liteClient.QueryAccount(address); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(v))
	}
}

func TestSendToken(t *testing.T) {
	initClient()
	memo := "airdrop test"
	addrAmountMap := map[string]float64{
		"faa1pa55mh8wuvdjsje7q3rnljv39krg3aeq3eu4mz": 13.23,
		"faa1j3ufmgwe2cuumj7423jt4creqlcskltn6ht5w9": 23.23,
		"faa1mqnv2t9cj4sps9ehmltvxwl24j5z8pt4mp4zmh": 313.23,
		"faa1q634ucwp92px3d8c9ayv8q9u4yduyfugkupe4c": 33.23,
		"faa1emzk4svf5fx8jc0f35krxv0u9hwtmvxhmzzv92": 213.23,
		"faa1p9x58hqtdclw9nsp20tnxc2nsrj9rl56y77yfs": 1323.23,
	}

	if len(addrAmountMap) == 0 {
		fmt.Println("data is empty")
		return
	}

	for k, v := range addrAmountMap {
		amount := util.Float64ToStr(v * math.Pow10(18))
		coin := ctypes.Coin{
			Denom:  "iris-atto",
			Amount: amount,
		}
		coins := []ctypes.Coin{coin}
		if res, err := txClient.SendToken(k, coins, memo, true); err != nil {
			fmt.Printf("%s fail, err: %s\n", k, err.Error())
		} else {
			if res.CommitResult.CheckTx.IsErr() || res.CommitResult.DeliverTx.IsErr() {
				fmt.Printf("%s fail, res: %s\n", k, util.ToJsonIgnoreErr(res))
			} else {
				fmt.Printf("%s:%v success, txHash: %s\n", k, v, res.CommitResult.Hash)
				fmt.Println("now sleep 10s")
				time.Sleep(time.Duration(10) * time.Second)
			}
		}
	}
}
