package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"nomadcoin/utils"
	"os"
)

func Start() {

}

type wallet struct {
	privateKey *ecdsa.PrivateKey
}

const walletFilename string = "nomadcoin.wallet"

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(walletFilename)
	return !os.IsNotExist(err) //file이 존재하면
}

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)
	return privKey
}
func persistKey(key *ecdsa.PrivateKey) {
	keybytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleError(err)
	err = os.WriteFile(walletFilename, keybytes, 0644)
	utils.HandleError(err)
}
func restoreKey() *ecdsa.PrivateKey {
	walletbytes, err := os.ReadFile(walletFilename)
	utils.HandleError(err)

	privKey, err := x509.ParseECPrivateKey(walletbytes)
	utils.HandleError(err)
	fmt.Println(privKey)
	return privKey
}

func Wallet() *wallet {
	if w == nil {
		var key *ecdsa.PrivateKey
		w = &wallet{}
		if hasWalletFile() {
			// yes -> restore from file
			key = restoreKey()
		} else {
			key = createPrivKey()
			persistKey(key)
		}
		w.privateKey = key
		//no -> create prv key, save to file
	}
	return w
}
