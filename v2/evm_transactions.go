package v2

import (
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/imskyd/mpc_generator/base"
)

// TokenTransfer basic token transfer
// if fee is nil, use cobo calc fee param
func (m *EvmMpcV2) TokenTransfer(from, to, tokenId, amount string, fee *coboWaas2.TransactionRequestFee) (*coboWaas2.CreateTransferTransaction201Response, error) {
	//source
	mpcSource := coboWaas2.NewMpcTransferSource(coboWaas2.WALLETSUBTYPE_ORG_CONTROLLED, m.walletId)
	mpcSource.SetAddress(from)
	source := coboWaas2.MpcTransferSourceAsTransferSource(mpcSource)
	//des
	des := coboWaas2.AddressTransferDestinationAsTransferDestination(coboWaas2.NewAddressTransferDestination(coboWaas2.TRANSFERDESTINATIONTYPE_ADDRESS))
	output := coboWaas2.NewAddressTransferDestinationAccountOutput(to, amount)
	des.AddressTransferDestination.AccountOutput = output
	//param
	transferParams := *coboWaas2.NewTransferParams(
		m.createRequestId(),
		source,
		tokenId,
		des,
	)
	if fee != nil {
		transferParams.Fee = fee
	}

	resp, _, err := m.client.TransactionsAPI.CreateTransferTransaction(m.getCtx()).TransferParams(transferParams).Execute()
	return resp, err
}

// TokenApprove basic contract call
// if fee is nil, use cobo calc fee param
func (m *EvmMpcV2) TokenApprove(chainId, from, token, spender, approveAmount string, fee *coboWaas2.TransactionRequestFee) (*coboWaas2.CreateTransferTransaction201Response, error) {
	calldata, _ := base.GetApproveCallData(spender, approveAmount)
	resp, err := m.ContractCall(chainId, from, token, calldata, "", fee)
	return resp, err
}

// ContractCall basic contract call
// if fee is nil, use cobo calc fee param
func (m *EvmMpcV2) ContractCall(chainId, from, to, callData, value string, fee *coboWaas2.TransactionRequestFee) (*coboWaas2.CreateTransferTransaction201Response, error) {
	//source
	mpcSource := coboWaas2.NewMpcContractCallSource(
		coboWaas2.CONTRACTCALLSOURCETYPE_ORG_CONTROLLED,
		m.walletId,
		from,
	)
	source := coboWaas2.MpcContractCallSourceAsContractCallSource(mpcSource)
	//des
	evmContractCallDes := coboWaas2.NewEvmContractCallDestination(coboWaas2.CONTRACTCALLDESTINATIONTYPE_EVM_CONTRACT, to, callData)
	if value != "" && value != "0" {
		evmContractCallDes.SetValue(value)
	}
	des := coboWaas2.EvmContractCallDestinationAsContractCallDestination(evmContractCallDes)
	//param
	contractCallParams := *coboWaas2.NewContractCallParams(
		m.createRequestId(),
		chainId,
		source,
		des,
	)
	if fee != nil {
		contractCallParams.Fee = fee
	}

	resp, _, err := m.client.TransactionsAPI.CreateContractCallTransaction(m.getCtx()).ContractCallParams(contractCallParams).Execute()
	return resp, err
}

/*
CancelTransaction
A transaction can be cancelled if its status is either of the following:
  - Submitted
  - PendingScreening
  - PendingAuthorization
  - PendingSignature
*/
func (m *EvmMpcV2) CancelTransaction(transactionId string) (*coboWaas2.CreateTransferTransaction201Response, error) {
	resp, _, err := m.client.TransactionsAPI.CancelTransactionById(m.getCtx(), transactionId).Execute()
	return resp, err
}

/*
DropTransaction
Dropping a transaction will trigger an Replace-By-Fee (RBF) transaction which is a new version of the original transaction. It has a higher transaction fee to incentivize miners to prioritize its confirmation over the original one. A transaction can be dropped if its status is Broadcasting.
  - For EVM chains, this RBF transaction has a transfer amount of 0 and the sending address is the same as the receiving address.
  - For UTXO chains, this RBF transaction has a transfer amount of 0 and the destination address is the same as the change address in the original transaction.
*/
func (m *EvmMpcV2) DropTransaction(transactionId string) (*coboWaas2.CreateTransferTransaction201Response, error) {
	resp, _, err := m.client.TransactionsAPI.DropTransactionById(m.getCtx(), transactionId).Execute()
	return resp, err
}
