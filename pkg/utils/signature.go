package utils

import "github.com/ethereum/go-ethereum/crypto"

func VerifySignature(pubdata []byte, msg []byte, sig []byte) bool {
	return crypto.VerifySignature(pubdata, msg, sig[:len(sig)-1])
}
