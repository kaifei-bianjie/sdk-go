package ad

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go"
	sdktypes "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/types/store"
	"github.com/irisnet/sdk-go/util"
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
	fromPassword = "eKF3TOm0BX"
	gasLimit     = uint64(150000)
	denom        = "uiris"
	feeAmount    = 30000
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
	// prod
	nodeUri := "http://seed-1.mainnet.irisnet.org:26657"
	grpcAddr := "seed-1.mainnet.irisnet.org:9090"
	chainId := "irishub-1"
	initKMType := "ks"

	// dev
	//network := types.Testnet
	//nodeUri := "http://192.168.150.40:26657"
	//grpcAddr := "192.168.150.40:9090"
	//initKMType := "ks"
	//chainId := "bifrost-2"

	// qa
	//nodeUri := "http://192.168.150.60:26657"
	//grpcAddr := "192.168.150.60:29090"
	//chainId := "irishub-qa2"
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
