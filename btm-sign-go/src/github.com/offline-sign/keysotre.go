package offline_sign

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bytom/common"
	"github.com/bytom/consensus"
	"github.com/bytom/crypto"
	"github.com/bytom/crypto/ed25519/chainkd"
	"github.com/bytom/protocol/bc"
	"github.com/bytom/protocol/vm/vmutil"
)




//生成私钥，公钥
func GenerateKey()(chainkd.XPrv, chainkd.XPub){
	//生成公钥私钥
	xpriv,xpub,err:=chainkd.NewXKeys(nil)
	if err!=nil {
		fmt.Errorf("create priv_key error,please try again,Err= %v",err)
		panic(err)
	}
	return xpriv,xpub
}
//公钥转换为地址
func XpubBytesToAddress(pub []byte)string{
	xpub:=BytesToXPub(pub)
	address:=XpubToAddress(xpub)
	return address
}

func XpubToAddress(xpub chainkd.XPub)string{
	pub := xpub.PublicKey()
	pubHash := crypto.Ripemd160(pub)

	//  TODO 切换生成主网还是测试网
	address, err := common.NewAddressWitnessPubKeyHash(pubHash, &consensus.TestNetParams)		//测试网
	//address, err := common.NewAddressWitnessPubKeyHash(pubHash, &consensus.ActiveNetParams)		//主网
	if err != nil {
		fmt.Errorf("create address error,please try again")
		panic(err)
	}
	return address.EncodeAddress()
}

//公钥转换为program
func XpubToProgramByte(xpub chainkd.XPub)[]byte{
	pub := xpub.PublicKey()
	pubHash := crypto.Ripemd160(pub)
	program, err := vmutil.P2WPKHProgram([]byte(pubHash))
	if err != nil {
		fmt.Errorf("create program error,please try again")
		panic(err)
	}
	return program
}

func XpubToProgramString(xpub chainkd.XPub)string{
	pub := xpub.PublicKey()
	pubHash := crypto.Ripemd160(pub)
	program, err := vmutil.P2WPKHProgram([]byte(pubHash))
	if err != nil {
		fmt.Println("create program error,please try again")
		panic(err)
	}
	return hex.EncodeToString(program)
}

//验证地址是否可用
func ValidAddress(address string)bool{
	_,err:=AddressToProgram(address)
	if err != nil {
		return false
	}
	return true
}

//将string类型的地址转换为program
func AddressToProgram(address string)([]byte,error){
	//TODO  先修改为测试网,可以更改未主网
	addr,err:=common.DecodeAddress(address,&consensus.TestNetParams)
	//addr,err:=common.DecodeAddress(address,&consensus.ActiveNetParams)
	if err != nil {
		return nil, err
	}
	redeemContract := addr.ScriptAddress()
	switch addr.(type) {
	case *common.AddressWitnessPubKeyHash:

		program,err:=vmutil.P2WPKHProgram(redeemContract)
		return program,err
	case *common.AddressWitnessScriptHash:
		program,err:=vmutil.P2WSHProgram(redeemContract)
		return program,err
	default:
		return nil,errors.New("Do not have this type address")
	}
}
//bytes 字节 转换为Xprv ----->
func BytesToXprv(res []byte)(xpriv chainkd.XPrv){
	copy(xpriv[:],res[:])
	return xpriv
}
func BytesToXPub(res []byte)(xpub chainkd.XPub){
	copy(xpub[:],res[:])
	return xpub
}

//string--->hash
//主要用于将string类型的Mux_id转换为bc.Hash
func MustDecodeHash(s string)(h bc.Hash,err error){
	if err := h.UnmarshalText([]byte(s)); err != nil {
		return bc.Hash{},err
	}
	return h,nil
}