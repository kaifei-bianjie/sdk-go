package ad

import (
	"encoding/json"
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
	network := types.Mainnet
	lcdUrl := "http://v2.irisnet-lcd.rainbow.one"
	rpcUrl := "tcp://seed-1.mainnet.irisnet.org:26657"
	initKMType := "ks"

	//network := types.Testnet
	//lcdUrl := "http://irisnet-lcd.dev.bianjie.ai"
	//rpcUrl := "tcp://192.168.150.32:26657"
	//initKMType := "ks"

	switch initKMType {
	case "seed":
		p := MnemonicKMParams{
			Menmonic: "proof own domain feature vehicle excite exotic way monkey stuff animal gorilla security roast street artwork room blue smoke fancy address",
			Password: "",
			FullPath: "44'/118'/0'/0/0",
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
			FilePath: "./ks.json",
			Password: "admino0o0oo0",
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
	address := "iaa16yued7d5m00nm9tu02pldxc8rgd57mvn5z8z0n"
	if v, err := liteClient.QueryAccount(address); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(v))
	}
}

func TestSendToken(t *testing.T) {
	initClient()
	memo := "airdrop from irisnet ama"
	addrAmountMap := map[string]float64{}

	//jsonDataStr := `{"iaa16zdqsdvca95x049qrdzvhmjn0hpveyv7hnul67": 0.2, "iaa1qtm7hw464wkr473k8qu30hgvh22gatdkc3w7p6": 0.1}`
	jsonDataStr := `{"iaa123uusha62lp7mgceg88dydaapjtt5p667xyez8":344.74,"iaa14l58vs85j3ql4f2ekzr2f9png9wsldwaugnaxj":344.74,"iaa17v7fn85prqpsff0ctlp7cr9s8zyg56vmaw87u8":344.74,"iaa17wnt8vuke6wu0gwxhq2nf8ej7nvpr6syh4nfs8":344.74,"iaa184p4x52677n6z66qnk2aatmrtnlfn0ls70d02d":344.74,"iaa18dgeshjgpdsknm7854296z3pmp6kekez5hnrl9":689.48,"iaa19v3dm8ypqzvd8y79y3vmswy93w9c76xfwx6evs":344.74,"iaa1affekf6ce3fcd0nhc4l30f5v32wxpalve99dlv":344.74,"iaa1axp0rxp449gmg23kqpxhhqk0lput8g27erdwdp":344.74,"iaa1c9d3srjxwp5t6yfjqhnj548glgvsxc2nqnlx4c":344.74,"iaa1cr79ek3frxy6hxkv5md52nkhxaqhygusscnzcn":344.74,"iaa1dm4cl6zrhspfc3x77d057mkagrnwpl8j2uyuug":344.74,"iaa1edp6gwutl0djtrvzwvvhl560p789l4eegz9sd2":344.74,"iaa1gah28g4cc8nuet5t3ndyw8606r8n226crgaxce":344.74,"iaa1gre3kkqee6x7ydstfqv5veddt7v88tsj5um34j":344.74,"iaa1gs8566mz04h86244dws24r05m3ufm0ycp0zhkx":861.85,"iaa1h2z2ful7jh02l2gclwlfq9zmhlw8yuuue3re0j":344.74,"iaa1hl6vfqwdcge8aasd503vze457fzaey8fal0rr4":344.74,"iaa1hvcugk8wk33laggy0a530dtdyln4k8wk6jtqhl":344.74,"iaa1j8efvwvvha8kdruwruwc0tgqu9aufkxzwukn2m":344.74,"iaa1jnpnym2h2h7ju3gy9pup78vf2j3krn9c8u0w2d":344.74,"iaa1lglrja7jvgc35g738zws565ews4s4vt0qxl408":344.74,"iaa1mqhjtzucnl3zgmdgddac3fjg9et463rcl40s24":861.85,"iaa1n2fqqq4xps5le5up268kny5hknzxfjl9l8tksf":344.74,"iaa1nma5007t8ne5a97mnlrwp4k0wje87k5atg8vpk":344.74,"iaa1ps8f37vr77pss2h2le7mzqsyetxmft53rvfnrr":344.74,"iaa1r5vsl2w3xac0puvruhh3v0c8ql37pxn9rpt4pv":344.74,"iaa1vew2r0zm46cuy5enk3sylecq778lkfd0gnaj2h":344.74,"iaa1vq6wf7llnrrrk594m6x2yhlrdlg5a0kjp9w57z":344.74,"iaa1vrhrgdlrugv6nxdlnjcnpd3q3xc98zg97f5erz":344.74,"iaa1vwzx5a0hgsreul076ev4wmux9xcqz58arunxsl":344.74,"iaa1ypjrytq487m2y9auqjct32ewsvg79v9l4uw5p0":344.74,"iaa1z2nn7qq2mz9vt9uq23t2afdmkdfkajzpg7gkfn":344.74,"iaa1zdcpsj0czadj5vw775v94xftz7dvg2j7efsqkm":344.74}`
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
		if acc, err := liteClient.QueryAccount(addr); err != nil {
			err := fmt.Errorf("get sender acc info fail, addr:%s, err:%s", addr, err.Error())
			return accNumber, sequence, err
		} else {
			if v, err := util.StrToInt64(acc.Value.AccountNumber); err != nil {
				err := fmt.Errorf("strToInt64 fail, str:%s. err:%s", acc.Value.AccountNumber, err.Error())
				return accNumber, sequence, err
			} else {
				accNumber = uint64(v)
			}

			if v, err := util.StrToInt64(acc.Value.Sequence); err != nil {
				err := fmt.Errorf("strToInt64 fail, str:%s. err:%s", acc.Value.Sequence, err.Error())
				return accNumber, sequence, err
			} else {
				sequence = uint64(v)
			}
		}

		return accNumber, sequence, nil
	}

	var (
		signerAccNumber uint64
		signerSequence  uint64
	)
	sender := km.GetAddr().String()
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
		amount := util.Float64ToStr(v * math.Pow10(18))
		coin := ctypes.Coin{
			Denom:  "iris-atto",
			Amount: amount,
		}
		coins := []ctypes.Coin{coin}

		if res, err := txClient.SendTokenWithSpecAccountInfo(k, coins, signerAccNumber, signerSequence, memo, false); err != nil {
			fmt.Printf("%s fail, err: %s\n", k, err.Error())
			time.Sleep(time.Duration(5) * time.Second)
		} else {
			if res.BroadcastResult.Code != 0 {
				fmt.Printf("%s fail, res: %s\n", k, util.ToJsonIgnoreErr(res))
				time.Sleep(time.Duration(10) * time.Second)
				if v1, v2, err := getAccInfo(sender); err != nil {
					fmt.Printf("query acc info fail, addr:%s, err:%s\n", sender, err.Error())
				} else {
					signerAccNumber = v1
					signerSequence = v2
				}
			} else {
				fmt.Printf("%s:%v success, txHash: %s\n", k, v, res.BroadcastResult.Hash)
				signerSequence += 1
			}
		}

		index += 1
	}
}
