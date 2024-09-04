package v2

import (
	"fmt"
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"log"
	"time"
)

func (m *EvmMpcV2) WaitTransactionDone(transactionId string) error {
	for {
		resp, err := m.GetTransactionByTransactionId(transactionId)
		if err != nil {
			return err
		}

		switch resp.Status {
		case coboWaas2.TRANSACTIONSTATUS_COMPLETED:
			//done
			return nil
		case coboWaas2.TRANSACTIONSTATUS_REJECTED:
			return fmt.Errorf("transaction rejected: tx id: %s", transactionId)
		case coboWaas2.TRANSACTIONSTATUS_FAILED:
			return fmt.Errorf("transaction failed: tx id: %s", transactionId)
		default:
			log.Printf("transaction id: %s status: %s", transactionId, resp.Status)
		}

		time.Sleep(3 * time.Second)
	}
}

func (m *EvmMpcV2) GetTransactionByTransactionId(transactionId string) (*coboWaas2.TransactionDetail, error) {
	resp, _, err := m.client.TransactionsAPI.GetTransactionById(m.getCtx(), transactionId).Execute()
	return resp, err
}
