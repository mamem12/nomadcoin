package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/nomadcoin/utils"
)

type wallet struct {
	PRIKey  *ecdsa.PrivateKey
	Address string
}

const (
	fileName string = "nomadcoin.wallet"
)

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(fileName)

	return !os.IsNotExist(err)

}

func createPRIKey() *ecdsa.PrivateKey {
	PRIKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return PRIKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(fileName, bytes, 0644)
	utils.HandleErr(err)
}

func restoreKey() *ecdsa.PrivateKey {
	keyAsBytes, err := os.ReadFile(fileName)
	utils.HandleErr(err)
	key, err := x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return key
}

func addressFromKey(key *ecdsa.PrivateKey) string {
	z := append(key.X.Bytes(), key.Y.Bytes()...)

	return fmt.Sprintf("%x", z)
}

func Wallet() *wallet {
	// has a wallet already?
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			// restore from file
			w.PRIKey = restoreKey()
		} else {
			key := createPRIKey()
			persistKey(key)
			w.PRIKey = key
		}
		w.Address = addressFromKey(w.PRIKey)
	}
	return w
}
