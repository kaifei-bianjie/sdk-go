package keys

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cosmos/go-bip39"
	"github.com/irisnet/irishub/crypto/keys/hd"
	"github.com/irisnet/irishub/crypto/keystore/uuid"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/types"
	ctypes "github.com/irisnet/irishub/types"
	"github.com/irisnet/sdk-go/types/tx"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"golang.org/x/crypto/pbkdf2"
	"io/ioutil"
)

type KeyManager interface {
	Sign(msg tx.StdSignMsg) ([]byte, error)
	GetPrivKey() crypto.PrivKey
	GetAddr() types.AccAddress
	ExportKeyStore(password string) (*EncryptedKeyJSON, error)
}

type keyManager struct {
	privKey  crypto.PrivKey
	addr     types.AccAddress
	mnemonic string
}

func (k *keyManager) Sign(msg tx.StdSignMsg) ([]byte, error) {
	sig, err := k.makeSignature(msg)
	if err != nil {
		return nil, err
	}

	newTx := auth.NewStdTx(msg.Msgs, msg.Fee, []auth.StdSignature{sig}, msg.Memo)
	bz, err := tx.Cdc.MarshalBinaryLengthPrefixed(newTx)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

func (k *keyManager) GetPrivKey() crypto.PrivKey {
	return k.privKey
}

func (k *keyManager) GetAddr() types.AccAddress {
	return k.addr
}

func (k *keyManager) ExportKeyStore(password string) (*EncryptedKeyJSON, error) {
	return generateKeyStore(k.privKey, password)
}

func (k *keyManager) makeSignature(msg tx.StdSignMsg) (sig auth.StdSignature, err error) {
	if err != nil {
		return
	}
	sigBytes, err := k.privKey.Sign(msg.Bytes())
	if err != nil {
		return
	}
	return auth.StdSignature{
		AccountNumber: msg.AccountNumber,
		Sequence:      msg.Sequence,
		PubKey:        k.privKey.PubKey(),
		Signature:     sigBytes,
	}, nil
}

func (k *keyManager) recoveryFromKeyStore(keystoreFile string, auth string) error {
	if auth == "" {
		return fmt.Errorf("Password is missing ")
	}
	keyJson, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		return err
	}
	var encryptedKey EncryptedKeyJSON
	err = json.Unmarshal(keyJson, &encryptedKey)
	if err != nil {
		return err
	}
	keyBytes, err := decryptKey(&encryptedKey, auth)
	if err != nil {
		return err
	}
	if len(keyBytes) != 32 {
		return fmt.Errorf("Len of Keybytes is not equal to 32 ")
	}
	var keyBytesArray [32]byte
	copy(keyBytesArray[:], keyBytes[:32])
	privKey := secp256k1.PrivKeySecp256k1(keyBytesArray)
	addr := ctypes.AccAddress(privKey.PubKey().Address())
	k.addr = addr
	k.privKey = privKey
	return nil
}

func (k *keyManager) recoverFromMnemonic(mnemonic, password, fullPath string) error {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		return err
	}

	masterPriv, chainCode := hd.ComputeMastersFromSeed(seed)
	privateKey, err := hd.DerivePrivateKeyForPath(masterPriv, chainCode, fullPath)

	if err != nil {
		return err
	}

	k.privKey = secp256k1.PrivKeySecp256k1(privateKey)
	k.addr = types.AccAddress(k.privKey.PubKey().Address())
	return nil

}

func generateKeyStore(privateKey crypto.PrivKey, password string) (*EncryptedKeyJSON, error) {
	addr := ctypes.AccAddress(privateKey.PubKey().Address())
	salt, err := GenerateRandomBytes(32)
	if err != nil {
		return nil, err
	}
	iv, err := GenerateRandomBytes(16)
	if err != nil {
		return nil, err
	}
	scryptParamsJSON := make(map[string]interface{}, 4)
	scryptParamsJSON["prf"] = "hmac-sha256"
	scryptParamsJSON["dklen"] = 32
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)
	scryptParamsJSON["c"] = 262144

	cipherParamsJSON := cipherparamsJSON{IV: hex.EncodeToString(iv)}
	derivedKey := pbkdf2.Key([]byte(password), salt, 262144, 32, sha256.New)
	encryptKey := derivedKey[:16]
	secpPrivateKey, ok := privateKey.(secp256k1.PrivKeySecp256k1)
	if !ok {
		return nil, fmt.Errorf(" Only PrivKeySecp256k1 key is supported ")
	}
	cipherText, err := aesCTRXOR(encryptKey, secpPrivateKey[:], iv)
	if err != nil {
		return nil, err
	}

	hasher := sha256.New()
	hasher.Write(derivedKey[16:32])
	hasher.Write(cipherText)
	mac := hasher.Sum(nil)

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	cryptoStruct := CryptoJSON{
		Cipher:       "aes-128-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          "pbkdf2",
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}
	return &EncryptedKeyJSON{
		Address: addr.String(),
		Crypto:  cryptoStruct,
		Id:      id.String(),
		Version: "1",
	}, nil
}

func NewKeyStoreKeyManager(file string, auth string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromKeyStore(file, auth)
	return &k, err
}

func NewMnemonicKeyManager(mnemonic, password, fullPath string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoverFromMnemonic(mnemonic, password, fullPath)
	return &k, err
}
