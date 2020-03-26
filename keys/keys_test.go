package keys

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/sdk-go/util"
	"github.com/irisnet/sdk-go/util/constant"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewKeyStoreKeyManager(t *testing.T) {
	file := "./ks_1234567890.json"
	if km, err := NewKeyStoreKeyManager(file, "1234567890"); err != nil {
		t.Fatal(err)
	} else {
		msg := []byte("hello world")
		signature, err := km.GetPrivKey().Sign(msg)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(km.GetAddr().String())

		assert.Equal(t, km.GetPrivKey().PubKey().VerifyBytes(msg, signature), true)
	}
}

func TestNewMnemonicKeyManager(t *testing.T) {
	sdk.SetNetworkType(constant.NetworkTypeMainnet)

	mnemonic := "proof general own domain feature brass pen vehicle excite exotic way monkey stuff animal gorilla security roast street artwork room blue smoke fancy address"
	password := ""
	fullpath := "44'/118'/0'/0/0"

	if km, err := NewMnemonicKeyManager(mnemonic, password, fullpath); err != nil {
		t.Fatal(err)
	} else {
		t.Log(km.GetAddr().String())

		if res, err := km.ExportKeyStore("secret"); err != nil {
			t.Fatal(err)
		} else {
			t.Log(util.ToJsonIgnoreErr(res))
		}
	}
}
