package offline_sign

import (
	"fmt"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	xpriv,xpub:=GenerateKey()
	fmt.Println(xpriv.String())
	fmt.Println(xpub.String())
	/*
	d0a28beaec724d103c6950ba6be2d4b010d19d19694834e32953eec085526a4231084b1b74a756f63b3f22e38ac7716ade2c9e6d0104725008baed62fb988abf
	524e5a065d0aa5897ffd59572d058e0835ba84e53930b9eb64755e8e4946b92731084b1b74a756f63b3f22e38ac7716ade2c9e6d0104725008baed62fb988abf
	*/
	address:=XpubToAddress(xpub)
	fmt.Println(address)
}


func TestOfflineSign(t *testing.T) {
	var (
		utxos []*UTXO
	)
	utxo:=UTXO{Amount:100000000,Address:"",MuxId:"",Position:0,CtrlProgram:""}
	utxos = append(utxos,&utxo)
	raw_tx,err:=BuildTransaction(utxos,"",90000000)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(raw_tx)
	//TODO submit raw_tx

}
