package offline_sign

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bytom/blockchain/txbuilder"
	"github.com/bytom/consensus"
	"github.com/bytom/consensus/segwit"
	"github.com/bytom/protocol/bc/types"
	"time"
)

type UTXO struct {
	OutId string
	MuxId string
	Position int
	Amount uint64
	CtrlProgram  string
	Address string
	Tag int
}

/*
根据UTXO  构建交易
*/
func BuildTransaction(utxos []*UTXO,wdAddress string,wdValue uint64)(string,error){
	var(
		ips []*types.TxInput
		ops []*types.TxOutput
		inputAmount uint64
	)
	//构建交易输入
	for _,utxo:=range utxos{
		value:=utxo.Amount
		//计算所有的input value
		inputAmount+=value
		input,err:=newBuildInput(utxo.MuxId,utxo.Address,utxo.CtrlProgram,value,uint64(utxo.Position))
		if err != nil {
			return "",err
		}
		ips = append(ips,input)
	}

	//构建tx
	tx:=types.NewTx(types.TxData{
		Version:1,
		Inputs:ips,
	})
	//计算手续费
	gas:=estimateTxGas(tx)

	//构建提币 output
	wdOuput,err:=newBuildOutput(wdAddress,wdValue)
	ops = append(ops,wdOuput)
	if err != nil {
		return "",err
	}
	//构建找零 output

	changeAddress:="tm1q2pz69vfkp3r4h7etsf8vs22eyrqrqalujh4r6e"  //找零地址
	if changeAddress==""{
		return "",errors.New("找零地址为空，无法找零")
	}
	//计算找零金额
	systemValue:=inputAmount-wdValue-uint64(gas)
	if systemValue<=0 {
		return "",errors.New(fmt.Sprintf("找零余额不足！！！，inputAmount=【%d】，wdValue=【%d】，gas=【%d】",inputAmount,wdValue,gas))
	}
	sysOutput,err:=newBuildOutput(changeAddress,systemValue)
	ops = append(ops,sysOutput)
	if err != nil {
		return "",errors.New(fmt.Sprintf("构建找零output错误，Err=【%v】",err))
	}
	//把 output 添加到tx中去
	tx.Outputs = ops

	//设置交易时间  这里设置为30分钟，如果30分钟没有上链，这笔交易就不在上链
	tx.TimeRange = uint64(time.Now().Unix()+1800)

	//构建template，
	tpl:=txbuilder.BuildTransaction(tx)  ///该方法写在与那么里面，当时为了方便直接继承它 的私有的接口了，也可以不写在源码里面
											//构建template时直接构造raw_tx类型的就可以。

	//交易输入的地址，
	// 所有的交易输入的地址都是一样的，所以取第一个交易输入的地址
	address:=utxos[0].Address
	if address==""{
		return "",errors.New("交易输入的地址为空")
	}
	raw_tx:=SignTransaction(address,tpl)
	return raw_tx,nil
}

func newBuildInput(MuxId,address,program string,amount uint64,position uint64)(*types.TxInput,error){
	var(
		input *types.TxInput
		err	error
	)
	if MuxId==""{
		err = errors.New(fmt.Sprintf("Do not contain mux_id ,mux_id= %s",MuxId))
		return nil,err
	}
	if (!ValidAddress(address)){
		err = errors.New(fmt.Sprintf("Input Address in not validated , address=%s, program = %s",address,program))
		return nil,err
	}
	//计算 mux_id
	mux_id,errs:=MustDecodeHash(MuxId)
	if errs != nil {
		err = errors.New(fmt.Sprintf("Mux_id calculate error,Err=【%v】",errs))
		return nil,err
	}
	//解码 Program
	controlProgram,errss:=hex.DecodeString(program)
	if errss != nil {
		err = errors.New(fmt.Sprintf("Decoding program error,Err=【%v】",errss))
		return nil,err
	}

	//build a new input
	input = types.NewSpendInput(nil,mux_id,*consensus.BTMAssetID,amount,position,controlProgram)
	return input,nil
}

func newBuildOutput(address string ,amount uint64)(*types.TxOutput,error){
	var(
		out *types.TxOutput
		err error
	)
	recvProg,errs:=AddressToProgram(address)
	if errs != nil {
		err = errors.New(fmt.Sprintf("Output address to program error,Err=【%v】",errs))
		return nil,err
	}
	out = types.NewTxOutput(*consensus.BTMAssetID, amount, recvProg)
	return out,nil
}

//计算交易手续费
func estimateTxGas(tx *types.Tx)int64{
	var  totalWitnessSize, totalP2WPKHGas, totalP2WSHGas, totalIssueGas int64
	baseSize := int64(176) // inputSize(112) + outputSize(64)
	baseP2WPKHSize := int64(98)
	baseP2WPKHGas := int64(1409)

	//遍历所有的input集合
	for _,input:=range tx.TxData.Inputs{
		switch input.InputType() {
		case types.SpendInputType:

			controlProgram := input.ControlProgram()
			//单地址类型
			if segwit.IsP2WPKHScript(controlProgram){
				totalWitnessSize += baseP2WPKHSize
				totalP2WPKHGas += baseP2WPKHGas
			}
		}
		//var witnessSize, gas int64
		//for i:=0;i<pos;i++{
		//	//模拟RawTxSigWitness所需要的手续费
		//	witnessSize += 33*int64(1) + 65*int64(1)
		//	gas += 1131*int64(1) + 72*int64(1) + 659+27
		//}
		//totalWitnessSize+=witnessSize
		//totalP2WPKHGas+=gas
	}

	flexibleGas := int64(0)
	if totalP2WPKHGas>0{
		flexibleGas += baseP2WPKHGas + (baseSize+baseP2WPKHSize)*consensus.StorageGasRate
	}

	//总共的存储手续费

	totalTxSizeGas := (int64(tx.TxData.SerializedSize)+totalWitnessSize)* consensus.StorageGasRate

	//总共的交易手续费
	totalGas := totalTxSizeGas + totalP2WPKHGas + totalP2WSHGas + totalIssueGas + flexibleGas
	return totalGas*consensus.VMGasRate
}


