package v2

import coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"

func New1559Fee(tokenId, maxFee, priorityFee, gasLimit string) coboWaas2.TransactionRequestFee {
	fee := coboWaas2.NewTransactionRequestEvmEip1559Fee(maxFee, priorityFee, coboWaas2.FEETYPE_EVM_EIP_1559, tokenId)
	if gasLimit != "" && gasLimit != "0" {
		fee.SetGasLimit(gasLimit)
	}
	paramFee := coboWaas2.TransactionRequestEvmEip1559FeeAsTransactionRequestFee(fee)
	return paramFee
}

func NewLegacyFee(tokenId, gasPrice, gasLimit string) coboWaas2.TransactionRequestFee {
	fee := coboWaas2.NewTransactionRequestEvmLegacyFee(gasPrice, coboWaas2.FEETYPE_EVM_LEGACY, tokenId)
	if gasLimit != "" && gasLimit != "0" {
		fee.SetGasLimit(gasLimit)
	}
	paramFee := coboWaas2.TransactionRequestEvmLegacyFeeAsTransactionRequestFee(fee)
	return paramFee
}
