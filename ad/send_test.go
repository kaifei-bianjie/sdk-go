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
	fromPassword = "test"
	gasLimit     = uint64(100000)
	denom        = "uiris"
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
	nodeUri := "http://v2.irisnet-lcd.rainbow.one"
	grpcAddr := "tcp://seed-1.mainnet.irisnet.org:26657"
	chainId := ""
	initKMType := "ks"

	//network := types.Testnet
	//nodeUri := "http://irisnet-lcd.dev.bianjie.ai"
	//grpcAddr := "tcp://192.168.150.32:26657"
	//initKMType := "ks"

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

func TestInitClient(m *testing.T) {
	initClient()
}

func TestGetAccountInfo(t *testing.T) {
	initClient()
	address := "iaa1xzw0rvmhhhzzt60m0rv84pd3frjzpld3l2qpjh"
	if v, err := irisClient.QueryAccount(address); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(v))
	}
}

func TestSendToken(t *testing.T) {
	initClient()
	memo := "sim upgrade validators rewards"
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
	jsonDataStr := `{"iaa122n4gkawquc446c36p3g83ulgy5lj0d3msty8q":2000,"iaa139qzyvac90z4j5xk7nqy8t2jmrlt9qnsl3m27f":2000,"iaa13hzqsv4e28syzhal7xwqm44saljs05l6rxxzkh":2000,"iaa13r5w8ts0370curm095uk57n89z26uen0fgupr4":2000,"iaa13tut6l4d63c9zpyflt25uw4dtdpdqm9afs3cml":2000,"iaa13y3fd2ej9aa3s9sfx93kkc6p8vghwr089yqgj7":2000,"iaa1543nj4z07vjqztvu3358fr2z2hcp0qtmdghutn":2000,"iaa15gyr8x527htlnq7nxm3v9kj82a2pw3jlkrl8cr":2000,"iaa15pzhk2scgmw6w97ftg20nw0ltyguyx88cjwmsy":2000,"iaa16plp8cmfkjssp222taq6pv6mkm8c5pa92fuyk6":2000,"iaa16v7g4a790t6qrrxr4scl5s652a6ha3g87v9drl":2000,"iaa186cqesc4m9l9lghucr07cum4tpccuu0a78k77d":2000,"iaa18clrjdxeg3ule5aprt85mt7h4q3zj0ecx05039":2000,"iaa19prg9sga9vqqhn65w62ny9mqknazl3lq53ffpq":2000,"iaa19wztuq7klrj2hgl4cqyvfajc6qnf24adtux2kz":2000,"iaa1aj6frfd5g0a5d88j6mjpvy0ql5afd5cr759nap":2000,"iaa1axd5c473vh6df4r0jevtzudhfk6k09s49eg3f5":2000,"iaa1dcqv69vqsdcc5lnmup6vkpv370fmqxgd56y0pe":2000,"iaa1f5e7uq0zk3f487vn52vjqqwm7xude54uad488j":2000,"iaa1g9ydd2wdgsszq3y60d5meq8mfrz9u6m9hlrugv":2000,"iaa1gvq24a5vn7twjupf3l2t7pnd9l4fm7uwwm4ujp":2000,"iaa1jmf7n6rpnqfse2xakdumw2mn5c2s9gmk3pe5s5":2000,"iaa1k4l7jnkfj6nfzxd3su94s0kauhcsx2gz2ymlum":2000,"iaa1kfhee2nqrg64krqa97q3ufw9d0phzp3j83mhg4":2000,"iaa1mjqef3jkgksk59rtnz3ljz94easln6cm9rj5th":2000,"iaa1mm4v28m7sfw4emm05d9sk8d3l9al3ffj3dgqez":2000,"iaa1qq6cly56h0yd8tr5fq6vw94pvp4d44nhfjfy2t":2000,"iaa1qq93sapmdcx36uz64vvw5gzuevtxsc7ldcvlqv":2000,"iaa1qtq8dwpdth5nwmyw60rm4texdnznk9ldsunptr":2000,"iaa1tp3rtyl8fhat9tkv5ufspz2zpdgdgunsgkjvgx":2000,"iaa1tsjrct9p7z2znsu4ehs69w5ydu5d5mu49h6yrk":2000,"iaa1ugjh8pn5sgvc8w45qgqrdkq4y74staqylzcwdk":2000,"iaa1w2dakpuvh9mglcs54wayta5dyv8vj8533tcdc7":2000,"iaa1z8wrnv35mmezpseym0jy7lngvsan2alwx9g2l5":2000,"iaa1zufrul2hvesqx5hrarvfcw34qnlq9yl8f6njd7":2000}`
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
		amt, err := util.StrToFloat64(util.Float64ToStrWithDecimal(v, 4))
		if err != nil {
			fmt.Printf("%d-%s fail, parse str to float64 err: %s\n", index, k, err.Error())
			continue
		}
		amountStr := util.Float64ToStr(amt * math.Pow10(6))
		amtInt, ok := sdktypes.NewIntFromString(amountStr)
		if !ok {
			fmt.Printf("%d-%s fail, parse str to amtInt err: %s\n", index, k, err.Error())
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
