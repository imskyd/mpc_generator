package mpc_generator

import (
	"encoding/json"
	"fmt"
	"github.com/CoboGlobal/cobo-go-api/cobo_custody"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
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

func EnvDev() cobo_custody.Env {
	return cobo_custody.Dev()
}

func EnvProd() cobo_custody.Env {
	return cobo_custody.Prod()
}

func (m *MPC) createRequestId() string {
	return fmt.Sprintf("cs-go-%d", time.Now().UnixMilli())
}

func wrapCoboErr(coboErr *cobo_custody.ApiError) error {
	if coboErr == nil {
		return nil
	}
	j, _ := json.Marshal(coboErr)
	return fmt.Errorf(string(j))
}

func (m *MPC) Transfer(coin, from, to string, amount, gasPrice, gasLimit *big.Int) (string, error) {
	requestId := m.createRequestId()
	res, err := m.client.CreateTransaction(coin, requestId, amount, from, to, "", nil, gasPrice, gasLimit, OperationTransfer, "", nil, nil, nil, "")
	if err != nil {
		return "", wrapCoboErr(err)
	}
	coboId, _ := res.Get("cobo_id").String()
	return coboId, nil
}

func (m *MPC) Approve(coin, from, token, spender string, approveAmount, gasPrice, gasLimit *big.Int) (string, error) {
	abiReader, _ := abi.JSON(strings.NewReader(ApproveABI))
	approvePack, packErr := abiReader.Pack("approve", common.HexToAddress(spender), approveAmount)
	if packErr != nil {
		return "", packErr
	}
	res, err := m.ContractCall(coin, from, token, fmt.Sprintf("0x%s", common.Bytes2Hex(approvePack)), big.NewInt(0), gasPrice, gasLimit)
	return res, err
}

func (m *MPC) ContractCall(coin, from, to, callData string, amount, gasPrice, gasLimit *big.Int) (string, error) {
	requestId := m.createRequestId()
	extraParameters := fmt.Sprintf("{\"calldata\": \"%s\"}", callData)
	res, err := m.client.CreateTransaction(coin, requestId, amount, from, to, "", nil, gasPrice, gasLimit, OperationContractCall, extraParameters, nil, nil, nil, "")
	if err != nil {
		return "", wrapCoboErr(err)
	}
	coboId, _ := res.Get("cobo_id").String()
	return coboId, nil
}

func (m *MPC) Approve1559(coin, from, token, spender string, approveAmount, maxFee, maxPriorityFee, gasLimit *big.Int) (string, error) {
	abiReader, _ := abi.JSON(strings.NewReader(ApproveABI))
	approvePack, packErr := abiReader.Pack("approve", common.HexToAddress(spender), approveAmount)
	if packErr != nil {
		return "", packErr
	}
	res, err := m.ContractCall1559(coin, from, token, fmt.Sprintf("0x%s", common.Bytes2Hex(approvePack)), big.NewInt(0), maxFee, maxPriorityFee, gasLimit)
	return res, err
}

func (m *MPC) ContractCall1559(coin, from, to, callData string, amount, maxFee, maxPriorityFee, gasLimit *big.Int) (string, error) {
	requestId := m.createRequestId()
	extraParameters := fmt.Sprintf("{\"calldata\": \"%s\"}", callData)
	res, err := m.client.CreateTransaction(coin, requestId, amount, from, to, "", nil, nil, gasLimit, OperationContractCall, extraParameters, maxFee, maxPriorityFee, nil, "")
	if err != nil {
		return "", wrapCoboErr(err)
	}
	coboId, _ := res.Get("cobo_id").String()
	return coboId, nil
}

func (m *MPC) GetTransactionByRequestIds(requestIds string) ([]Transaction, error) {
	res, err := m.client.TransactionsByRequestIds(requestIds, 0)
	if err != nil {
		return nil, fmt.Errorf(err.ErrorMessage)
	}
	txJson, err2 := res.Get("transactions").MarshalJSON()
	if err != nil {
		return nil, err2
	}
	var transactions []Transaction
	_ = json.Unmarshal(txJson, &transactions)
	return transactions, err2
}

func (m *MPC) GetTransactionByCoboIds(coboIds string) ([]Transaction, error) {
	res, err := m.client.TransactionsByCoboIds(coboIds, 0)
	if err != nil {
		return nil, fmt.Errorf(err.ErrorMessage)
	}
	txJson, err2 := res.Get("transactions").MarshalJSON()
	if err != nil {
		return nil, err2
	}
	var transactions []Transaction
	_ = json.Unmarshal(txJson, &transactions)
	return transactions, err2
}

func (m *MPC) GetTransactionsListByAddress(address string) ([]Transaction, error) {
	res, err := m.client.ListTransactions(0, 0, 0, "", "", 0,
		"", address, "", 0)
	if err != nil {
		return nil, fmt.Errorf(err.ErrorMessage)
	}
	txJson, err2 := res.Get("transactions").MarshalJSON()
	if err != nil {
		return nil, err2
	}
	var transactions []Transaction
	_ = json.Unmarshal(txJson, &transactions)
	return transactions, err2
}

func GetTxStatusContent(status int) string {
	if val, ok := TextTxStatus[status]; ok {
		return val
	}
	return fmt.Sprintf("Undefined:%d", status)
}
