package v2

import (
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

/*
ResendTransaction
This operation resends a specified transaction.

	Resending a transaction initiates a new attempt to process the transaction that failed previously.
	A transaction can be resent if its status is failed.
*/
func (m *EvmMpcV2) ResendTransaction(transactionId string) (*coboWaas2.CreateTransferTransaction201Response, error) {
	transactionResend := *coboWaas2.NewTransactionResend(m.createRequestId()) // TransactionResend | The request body to resend transactions (optional)
	resp, r, err := m.client.TransactionsAPI.ResendTransactionById(m.getCtx(), transactionId).TransactionResend(transactionResend).Execute()
	return m.formatResponse(resp, r, err)
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
	resp, r, err := m.client.TransactionsAPI.CancelTransactionById(m.getCtx(), transactionId).Execute()
	return m.formatResponse(resp, r, err)
}

/*
DropTransaction
Dropping a transaction will trigger an Replace-By-Fee (RBF) transaction which is a new version of the original transaction. It has a higher transaction fee to incentivize miners to prioritize its confirmation over the original one. A transaction can be dropped if its status is Broadcasting.
  - For EVM chains, this RBF transaction has a transfer amount of 0 and the sending address is the same as the receiving address.
  - For UTXO chains, this RBF transaction has a transfer amount of 0 and the destination address is the same as the change address in the original transaction.
*/
func (m *EvmMpcV2) DropTransaction(transactionId string, fee coboWaas2.TransactionRequestFee) (*coboWaas2.CreateTransferTransaction201Response, error) {
	transactionRbf := coboWaas2.NewTransactionRbf(m.createRequestId(), fee)

	resp, r, err := m.client.TransactionsAPI.DropTransactionById(m.getCtx(), transactionId).TransactionRbf(*transactionRbf).Execute()
	return m.formatResponse(resp, r, err)
}

/*
SpeedUpTransaction
Speeding up a transaction will trigger an Replace-By-Fee (RBF) transaction which is a new version of the original transaction.

	    It shares the same inputs but has a higher transaction fee to incentivize miners to prioritize its confirmation over the previous one.
		A transaction can be accelerated if its status is Broadcasting.
*/
func (m *EvmMpcV2) SpeedUpTransaction(transactionId string, fee coboWaas2.TransactionRequestFee) (*coboWaas2.CreateTransferTransaction201Response, error) {
	transactionRbf := *coboWaas2.NewTransactionRbf(m.createRequestId(), fee) // TransactionRbf | The request body to drop or to speed up transactions (optional)

	resp, r, err := m.client.TransactionsAPI.SpeedupTransactionById(m.getCtx(), transactionId).TransactionRbf(transactionRbf).Execute()
	return m.formatResponse(resp, r, err)
}
