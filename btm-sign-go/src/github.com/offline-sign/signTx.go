package offline_sign

import (

	"encoding/hex"

	"fmt"
	"github.com/bytom/blockchain/txbuilder"
	"github.com/bytom/crypto/ed25519/chainkd"

)

func SignTransaction(address string,tpl *txbuilder.Template)string{

	if tpl.SigningInstructions==nil{
		tpl.SigningInstructions = []*txbuilder.SigningInstruction{}
	}
	var (
		newTpl *txbuilder.Template
	)
	for i,_:=range tpl.Transaction.Inputs{
		h:=tpl.Hash(uint32(i)).Byte32()
		//进行签名
		sig,xPub,_:=sign(address,h[:])
		data:=[]byte(xPub.PublicKey())
		fmt.Printf("签名数据[%d]：%s",i,hex.EncodeToString(sig))
		newTpl=txbuilder.BuildSignatureDataToTplByJun(tpl,sig,data,i)  //该方法写在与那么里面
	}
	tt,err:=txbuilder.CheckTpl(newTpl)
	if err != nil {
		return ""
	}
	fmt.Println(" Sign a transaction success !!!")

	return tt.Transaction.String()
}

func sign(address string,data []byte)(sig[]byte,pub chainkd.XPub,err error){
	var xPriv *chainkd.XPrv
	//todo 根据地址获取对应的私钥
	// xPriv:= todo
	if err != nil {
		return nil,chainkd.XPub{},err
	}
	//签名交易
	sig =xPriv.Sign(data[:])
	return sig,xPriv.XPub(),nil
}