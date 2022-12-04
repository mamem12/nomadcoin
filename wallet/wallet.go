package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
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

func Sing(payload string, w *wallet) string {
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	r, s, err := ecdsa.Sign(rand.Reader, nil, payloadAsBytes)
	utils.HandleErr(err)

	signature := append(r.Bytes(), s.Bytes()...)

	return fmt.Sprintf("%x", signature)
}

func restoreBigInt(payload string) (*big.Int, *big.Int, error) {
	signBytes, err := hex.DecodeString(payload)

	if err != nil {
		return nil, nil, err
	}

	firstHalfBytes := signBytes[:len(signBytes)/2]
	secondHalfBytes := signBytes[len(signBytes)/2:]

	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)

	return &bigA, &bigB, nil
}

func Verify(signature, payload, address string) bool {
	r, s, err := restoreBigInt(signature)
	utils.HandleErr(err)

	x, y, err := restoreBigInt(address)
	utils.HandleErr(err)

	PUBKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	return ecdsa.Verify(&PUBKey, payloadAsBytes, r, s)
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
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
