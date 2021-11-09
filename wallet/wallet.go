package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"nomadcoin/utils"
	"os"
)

func Start() {

}

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string //public key to hexadecimal
}

const walletFilename string = "nomadcoin.wallet"

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(walletFilename)
	return !os.IsNotExist(err) //file이 존재하면
}

func createPrivKey() (key *ecdsa.PrivateKey) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)
	return
}

func addFromKey(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}
func persistKey(key *ecdsa.PrivateKey) {
	keybytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleError(err)
	err = os.WriteFile(walletFilename, keybytes, 0644)
	utils.HandleError(err)
}
func restoreKey() (key *ecdsa.PrivateKey) {
	walletbytes, err := os.ReadFile(walletFilename)
	utils.HandleError(err)

	key, err = x509.ParseECPrivateKey(walletbytes)
	utils.HandleError(err)
	return
}
func encodeBigInts(a, b []byte) string {
	z := append(a, b...)

	//Privatekey의 publickey의 x,y를 합친게 address
	return fmt.Sprintf("%x", z)
}
func sign(payload string, w wallet) string {
	payloadAsB, err := hex.DecodeString(payload)
	utils.HandleError(err)

	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsB)
	utils.HandleError(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	bigA, bigB := big.Int{}, big.Int{}

	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)

	return &bigA, &bigB, nil
}

func verify(payload, signature string, address string) bool {
	r, s, err := restoreBigInts(signature)
	utils.HandleError(err)

	x, y, err := restoreBigInts(address) //address를 통해 publickey를 가져올 수 있음
	utils.HandleError(err)
	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleError(err)
	return ecdsa.Verify(publicKey, payloadBytes, r, s)
}
func Wallet() *wallet {
	if w == nil {
		var key *ecdsa.PrivateKey
		w = &wallet{}
		if hasWalletFile() {
			// yes -> restore from file
			w.privateKey = restoreKey()
		} else {
			key = createPrivKey()
			persistKey(key)
		}
		//no -> create prv key, save to file
	}
	w.Address = addFromKey(w.privateKey)
	return w
}
