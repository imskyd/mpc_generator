package v2

import (
	"fmt"
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"time"
)

func (m *EvmMpcV2) WaitTransactionDone(transactionId string, maxTryTime int) error {
	tryTime := 1
	for {
		if tryTime > maxTryTime {
			return fmt.Errorf("wait transaction get max try times")
		}
		resp, err := m.GetTransactionByTransactionId(transactionId)
		if err != nil {
			return err
		}

		switch resp.Status {
		case coboWaas2.TRANSACTIONSTATUS_COMPLETED:
			//done
			return nil
		case coboWaas2.TRANSACTIONSTATUS_REJECTED:
			return fmt.Errorf("transaction id: %s, request id: %s, cobo id: %s", transactionId, *resp.RequestId, *resp.CoboId)
		case coboWaas2.TRANSACTIONSTATUS_FAILED:
			return fmt.Errorf("transaction id: %s, request id: %s, cobo id: %s", transactionId, *resp.RequestId, *resp.CoboId)
		default:
			//continue when got other status
			m.logger.Infof("transaction id: %s, request id: %s, cobo id: %s, status: %s", transactionId, *resp.RequestId, *resp.CoboId, resp.Status)
		}

		time.Sleep(3 * time.Second)
		tryTime++
	}
}

func (m *EvmMpcV2) GetTransactionByTransactionId(transactionId string) (*coboWaas2.TransactionDetail, error) {
	resp, _, err := m.client.TransactionsAPI.GetTransactionById(m.getCtx(), transactionId).Execute()
	if err == nil && m.debug {
		m.printFormatLog("GetTransactionByTransactionId", resp)
	}
	return resp, err
}

/*
ListAllTransactions
  - request_id. The request ID that is used to track a transaction request. The request ID is provided by you and must be unique within your organization.
  - cobo_ids. A list of Cobo IDs, separated by comma. A Cobo ID can be used to track a transaction.
  - transaction_ids. A list of transaction IDs, separated by comma.
*/
func (m *EvmMpcV2) ListAllTransactions(requestId, transactionIds, coboIds string) (*coboWaas2.ListTransactions200Response, error) {
	request := m.client.TransactionsAPI.ListTransactions(m.getCtx())
	if requestId != "" {
		request = request.RequestId(requestId)
	}
	if transactionIds != "" {
		request = request.TransactionIds(transactionIds)
	}
	if coboIds != "" {
		request = request.CoboIds(coboIds)
	}
	resp, _, err := request.Execute()
	if err == nil && m.debug {
		m.printFormatLog("GetTransactionByTransactionId", resp)
	}
	return resp, err
}
