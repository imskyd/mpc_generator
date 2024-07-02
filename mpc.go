package mpc_generator

import (
	"fmt"
	"github.com/CoboGlobal/cobo-go-api/cobo_custody"
	"math/big"
	"time"
)

type MPC struct {
	client *cobo_custody.MPCClient
}

func NewMPC(apiKey, apiSecret string, env cobo_custody.Env) *MPC {
	fmt.Println("API_SECRET:", apiSecret)
	fmt.Println("API_KEY:", apiKey)

	var localSigner = cobo_custody.LocalSigner{
		PrivateKey: apiSecret,
	}

	var client = cobo_custody.MPCClient{
		Signer: localSigner,
		Env:    env,
	}

	var m MPC
	m.client = &client
	return &m
}

func (m *MPC) createRequestId() string {
	return fmt.Sprintf("cs-go-%d", time.Now().UnixMilli())
}

func (m *MPC) Transfer(coin, from, to string, amount, gasPrice, gasLimit *big.Int) (string, *cobo_custody.ApiError) {
	requestId := m.createRequestId()
	res, err := m.client.CreateTransaction(coin, requestId, amount, from, to, "", nil, gasPrice, gasLimit, OperationTransfer, "", nil, nil, nil, "")
	if err != nil {
		return "", err
	}
	coboId, _ := res.Get("cobo_id").String()
	return coboId, nil
}

func (m *MPC) ContractCall(coin, from, to, callData string, amount, gasPrice, gasLimit *big.Int) (string, *cobo_custody.ApiError) {
	requestId := m.createRequestId()
	extraParameters := fmt.Sprintf("{\"calldata\": \"%s\"}", callData)
	res, err := m.client.CreateTransaction(coin, requestId, amount, from, to, "", nil, gasPrice, gasLimit, OperationContractCall, extraParameters, nil, nil, nil, "")
	if err != nil {
		return "", err
	}
	coboId, _ := res.Get("cobo_id").String()
	return coboId, nil
}

func (m *MPC) GetTransactionByRequestIds(requestIds string) ([]Transaction, error) {
	res, err := m.client.TransactionsByRequestIds(requestIds, 0)
	if err != nil {
		return nil, fmt.Errorf(err.ErrorMessage)
	}
	array, err2 := res.Get("result").Get("transactions").Array()
	var transactions []Transaction
	for _, a := range array {
		transactions = append(transactions, a.(Transaction))
	}
	return transactions, err2
}

func GetTxStatusContent(status int) string {
	if val, ok := TextTxStatus[status]; ok {
		return val
	}
	return fmt.Sprintf("Undefined:%d", status)
}
