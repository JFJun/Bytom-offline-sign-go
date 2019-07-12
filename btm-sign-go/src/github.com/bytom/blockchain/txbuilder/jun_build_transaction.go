package txbuilder

import (
	chainjson "github.com/bytom/encoding/json"
	"github.com/bytom/protocol/bc/types"
)

//write by jun
//this is a method that to build a UTXO transaction
//and it will be sign a transaction


func BuildSignatureDataToTplByJun(tpl *Template,sig,data []byte,i int)( *Template){
	sigInst:=tpl.SigningInstructions[i]
	if len(sigInst.WitnessComponents)==0{
		witCom:=[]witnessComponent{
			&RawTxSigWitness{
				Quorum: 1,
				Sigs:   []chainjson.HexBytes{sig},
			},
			DataWitness(data),
		}
		sigInst.WitnessComponents=witCom
	}
	return tpl

}

func CheckTpl(tpl *Template)(*Template,error){
	if err := materializeWitnesses(tpl); err != nil {
		return nil,err
	}

	//if !testutil.DeepEqual(tx, tpl.Transaction) {
	//	return nil,errors.New(fmt.Sprintf("tx:%v result is equal to want:%v", tx, tpl.Transaction))
	//}
	return tpl,nil
}

func BuildTransaction(tx *types.Tx)*Template{
	tpl:=&Template{}
	tpl.AllowAdditional = false
	for i,_:=range tx.Inputs{
		instruction:=&SigningInstruction{}
		instruction.Position = uint32(i)
		// Empty signature arrays should be serialized as empty arrays, not null.
		if instruction.WitnessComponents ==nil{
			instruction.WitnessComponents = []witnessComponent{}
		}
		tpl.SigningInstructions = append(tpl.SigningInstructions,instruction)
	}
	tpl.Transaction = tx
	tpl.Fee = CalculateTxFee(tpl.Transaction)  //计算手续费
	return tpl
}