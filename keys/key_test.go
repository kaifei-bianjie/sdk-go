package keys

import (
	sdkcrypto "github.com/irisnet/irishub-sdk-go/crypto"
	sdktypes "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/sdk-go/util"
	"github.com/pkg/errors"
	"testing"
)

func TestNewKeyStoreKeyManager(t *testing.T) {
	file := "./ks.txt"
	passWd := "200200200"

	if c, err := util.ReadFile(file); err != nil {
		t.Fatalf("read ks file fail, err: %s", err.Error())
	} else {
		privKey, algo, err := sdkcrypto.UnarmorDecryptPrivKey(c, passWd)
		if err != nil {
			t.Fatal(errors.Wrap(err, "failed to decrypt private key"))
		}

		msg := []byte("hello world")
		signature, err := privKey.Sign(msg)
		if err != nil {
			t.Fatal(err)
		}

		if !privKey.PubKey().VerifySignature(msg, signature) {
			t.Fatal("VerifySignature fail")
		} else {
			t.Logf("privKey algo: %s, addr: %s", algo,
				sdktypes.AccAddress(privKey.PubKey().Address().Bytes()).String())
		}
	}
}

func TestNewMnemonicKeyManager(t *testing.T) {
	defaultAlgo := "secp256k1"
	mnemonic := ""
	hdPath := ""
	ksPasswd := ""

	if km, err := sdkcrypto.NewMnemonicKeyManagerWithHDPath(mnemonic, defaultAlgo, hdPath); err != nil {
		t.Fatal(err)
	} else {
		t.Log(sdktypes.AccAddress(km.ExportPubKey().Address().Bytes()).String())
		if res, err := km.ExportPrivKey(ksPasswd); err != nil {
			t.Fatal(err)
		} else {
			t.Log(util.ToJsonIgnoreErr(res))
		}
	}
}
