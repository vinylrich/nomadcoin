package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"nomadcoin/utils"
)

const (
	hashedmsg  string = "1c5863cd55b5a4413fd59f054af57ba3c75c0698b3851d70f99b8de2d5c7338f"
	privateKey string = "30770201010420f163e0eb74c179eb9843c08a19088734cd12ff20a534d1a1878f4e014f921a72a00a06082a8648ce3d030107a1440342000450d027f066e0e51087199fe65704041073ac92098f141c98e23932318ecd59c1af3d8da8018e88939db9873ef4dcc33876cd3a126385ad51db9c90212717c850"
	signature  string = "4229962a31e14127287112b8cbe4c5c8cac54c7e0dd21a286c25e1d334cd780875e0d904d39dffa1d5ccf4eef52d9ca2de02ee3af2ec358febebcdb9f5bf53f4"
)

func Start() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)
	keyAsBytes, err := x509.MarshalECPrivateKey(privateKey)
	utils.HandleError(err)

	fmt.Printf("%x\n\n", keyAsBytes)
	bytemsg, err := hex.DecodeString(hashedmsg)

	utils.HandleError(err)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, bytemsg)
	//signature

	signature := append(r.Bytes(), s.Bytes()...)

	fmt.Printf("%x\n", signature)
	utils.HandleError(err)
}
