package keypair

import (
	. "github.com/Lawliet-Chan/yu/common"
	"testing"
)

func TestSrKey(t *testing.T) {
	pubkey, privkey, err := GenKeyPair(Sr25519)
	if err != nil {
		panic("generate key error: " + err.Error())
	}
	ecall := &Ecall{
		TripodName: "asset",
		ExecName:   "Transfer",
		Params:     JsonString("params"),
	}

	signByt, err := privkey.SignData(ecall.Bytes())
	if err != nil {
		panic("sign data error: " + err.Error())
	}

	genPubkey, err := PubKeyFromBytes(pubkey.BytesWithType())
	if err != nil {
		panic("gen pubkey error: " + err.Error())
	}
	t.Logf("verify signature result:  %v", genPubkey.VerifySignature(ecall.Bytes(), signByt))
}
