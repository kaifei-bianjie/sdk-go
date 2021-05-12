package ad

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub-sdk-go"
	sdktypes "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/types/store"
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

const (
	fromName     = "test"
	fromPassword = "200200200"
	gasLimit     = uint64(100000)
	denom        = "ubif"
	feeAmount    = 20000
)

var (
	irisClient *sdk.IRISHUBClient
	sendAddr   string
	fee        = sdktypes.DecCoin{
		Denom:  denom,
		Amount: sdktypes.NewDec(feeAmount),
	}

	keyDaoOption = sdktypes.KeyDAOOption(store.NewMemory(nil))
)

func initClient() {
	//nodeUri := "http://sentry-0.mainnet.irisnet.org:26657"
	//grpcAddr := "sentry-1.mainnet.irisnet.org:9090"
	//chainId := "irishub-1"
	//initKMType := "ks"

	//network := types.Testnet
	nodeUri := "http://192.168.150.40:26657"
	grpcAddr := "192.168.150.40:9090"
	initKMType := "ks"
	chainId := "bifrost-2"

	options := []sdktypes.Option{keyDaoOption}
	if cfg, err := sdktypes.NewClientConfig(nodeUri, grpcAddr, chainId, options...); err != nil {
		panic(err)
	} else {
		c := sdk.NewIRISHUBClient(cfg)
		irisClient = &c
	}

	switch initKMType {
	case "seed":
		p := MnemonicKMParams{
			Menmonic: "",
			Password: fromPassword,
			FullPath: "44'/118'/0'/0/0",
		}
		if v, err := irisClient.Key.RecoverWithHDPath(fromName, p.Password, p.Menmonic, p.FullPath); err != nil {
			panic(fmt.Errorf("recover client key fail, err: %s", err.Error()))
		} else {
			sendAddr = v
		}
		break
	case "ks":
		p := KeyStoreKMParams{
			FilePath: "./ks.txt",
			Password: fromPassword,
		}
		content, err := util.ReadFile(p.FilePath)
		if err != nil {
			panic(fmt.Sprintf("read ks file fail, err: %s", err.Error()))
		}
		if v, err := irisClient.Key.Import(fromName, p.Password, content); err != nil {
			panic(fmt.Errorf("recover client key fail, err: %s", err.Error()))
		} else {
			sendAddr = v
		}
		break
	default:
		panic("should init km first")
	}
}

func TestInitClient(t *testing.T) {
	initClient()
	t.Log(sendAddr)
}

func TestGetAccountInfo(t *testing.T) {
	initClient()
	address := "iaa1tzl7vrq99l8r4yvnh3fl5k8y35q560fr09em8e"
	if v, err := irisClient.QueryAccount(address); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(v))
	}
}

func TestSendToken(t *testing.T) {
	initClient()
	toAddr := "iaa1czu8x6drhwry44sc3jtwu4mylc574l09x6zmh3"
	memo := "test"
	amount := float64(0.1)
	baseTx := sdktypes.BaseTx{
		From:     fromName,
		Password: fromPassword,
		Gas:      gasLimit,
		Fee:      sdktypes.DecCoins{fee},
		Memo:     memo,
		Mode:     sdktypes.Sync,
	}

	// query account info
	getAccInfo := func(addr string) (uint64, uint64, error) {
		var (
			accNumber uint64
			sequence  uint64
		)
		if acc, err := irisClient.QueryAccount(addr); err != nil {
			err := fmt.Errorf("get sender acc info fail, addr:%s, err:%s", addr, err.Error())
			return accNumber, sequence, err
		} else {
			accNumber = acc.AccountNumber
			sequence = acc.Sequence
		}

		return accNumber, sequence, nil
	}

	var (
		signerAccNumber uint64
		signerSequence  uint64
	)
	sender := sendAddr
	if v1, v2, err := getAccInfo(sender); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		signerAccNumber = v1
		signerSequence = v2
	}

	amountStr := util.Float64ToStrWithDecimal(amount*math.Pow10(6), 0)
	amtInt, ok := sdktypes.NewIntFromString(amountStr)
	if !ok {
		panic(fmt.Errorf(" parse str to amtInt fail, amount: %f\n", amount))
	}
	decCoin := sdktypes.NewDecCoin(denom, amtInt)
	decCoins := sdktypes.NewDecCoins(sdktypes.DecCoins{decCoin}...)
	if res, err := irisClient.Bank.SendWitchSpecAccountInfo(toAddr, signerSequence, signerAccNumber, decCoins, baseTx); err != nil {
		t.Fatalf("send token fail, errCode: %d, codeSpace: %s, err: %s", err.Code(), err.Codespace(), err.Error())
	} else {
		t.Logf("success, txHash: %s", res.Hash)
	}
}

func TestAirDrop(t *testing.T) {
	initClient()
	memo := "delegator rewards for 2021/02"
	baseTx := sdktypes.BaseTx{
		From:     fromName,
		Password: fromPassword,
		Gas:      gasLimit,
		Fee:      sdktypes.DecCoins{fee},
		Memo:     memo,
		Mode:     sdktypes.Sync,
	}

	addrAmountMap := map[string]float64{}

	//jsonDataStr := `{"iaa1dcmuvs0gqx88kmuncejt87x9q7pargk2zu5s8t": 0.01, "iaa16zdqsdvca95x049qrdzvhmjn0hpveyv7hnul67": 0.01}`
	jsonDataStr := `{"iaa1tuxtl5wmfavkhxq5cu9mvhw2xg622avpwex0yr":24.2538423850324,"iaa1urz36q72829jwk34yzpp2zyms6e977cf6a6jd5":14.2175736699461,"iaa12whfk6ycqv4ym2kl9a4rr5apyp7s5ks60yp05h":1.28704283080337,"iaa1dwzxt9v3qrd65ut5zj738yg4u8z496nhxz9vwc":216.386639284562}`
	if err := json.Unmarshal([]byte(jsonDataStr), &addrAmountMap); err != nil {
		t.Fatalf("unmarshal json fail, err is %s\n", err.Error())
	}

	handledDataStr := ``
	var handledAddrList []string
	if handledDataStr != "" {
		if err := json.Unmarshal([]byte(handledDataStr), &handledAddrList); err != nil {
			t.Fatalf("unmarshal json fail, err is %s\n", err.Error())
		}
	}

	if len(addrAmountMap) == 0 {
		fmt.Println("data is empty")
		return
	}

	listToMap := func(list []string) map[string]bool {
		res := make(map[string]bool)
		if len(list) > 0 {
			for _, v := range list {
				res[v] = true
			}
		}
		return res
	}
	handledAddrMap := listToMap(handledAddrList)

	// query account info
	getAccInfo := func(addr string) (uint64, uint64, error) {
		var (
			accNumber uint64
			sequence  uint64
		)
		if acc, err := irisClient.QueryAccount(addr); err != nil {
			err := fmt.Errorf("get sender acc info fail, addr:%s, err:%s", addr, err.Error())
			return accNumber, sequence, err
		} else {
			accNumber = acc.AccountNumber
			sequence = acc.Sequence
		}

		return accNumber, sequence, nil
	}

	var (
		signerAccNumber uint64
		signerSequence  uint64
	)
	sender := sendAddr
	if v1, v2, err := getAccInfo(sender); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		signerAccNumber = v1
		signerSequence = v2
	}

	index := 1
	for k, v := range addrAmountMap {
		if index%100 == 0 {
			time.Sleep(10 * time.Second)
		}
		if handledAddrMap[k] {
			fmt.Printf("%s has been handled\n", k)
			continue
		}
		amountStr := util.Float64ToStrWithDecimal(v*math.Pow10(6), 0)
		amtInt, ok := sdktypes.NewIntFromString(amountStr)
		if !ok {
			fmt.Printf("%d-%s fail, parse str to amtInt, amt: %.8f, amtStr: %s\n",
				index, k, v, amountStr)
			continue
		}

		decCoin := sdktypes.NewDecCoin(denom, amtInt)
		decCoins := sdktypes.NewDecCoins(sdktypes.DecCoins{decCoin}...)

		if res, err := irisClient.Bank.SendWitchSpecAccountInfo(k, signerSequence, signerAccNumber, decCoins, baseTx); err != nil {
			fmt.Printf("%d-%s fail, errCode: %d, codeSpace: %s, err: %s\n",
				index, k, err.Code(), err.Codespace(), err.Error())
			time.Sleep(time.Duration(5) * time.Second)

			if v1, v2, err := getAccInfo(sender); err != nil {
				fmt.Printf("query acc info fail, addr:%s, err:%s\n", sender, err.Error())
			} else {
				signerAccNumber = v1
				signerSequence = v2
			}
		} else {
			fmt.Printf("%d-%s:%v success, txHash: %s\n", index, k, v, res.Hash)
			signerSequence += 1
		}

		index += 1
	}
}
