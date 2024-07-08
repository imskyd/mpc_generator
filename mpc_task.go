package mpc_generator

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
	"strings"
	"time"
)

type Task struct {
	TaskId    string
	coin      string
	from      string
	to        string
	callData  string
	operation int
	amount    *big.Int
	gasPrice  *big.Int
	gasLimit  *big.Int
}

func (m *MPC) RunTaskFlow(tasks []*Task) {
	for true {
		for _, task := range tasks {
			for true {
				txs, err := m.GetTransactionsListByAddress(task.from)
				if err != nil {
					log.Printf("m.GetTransactionsListByAddress(%s) err: %s", task.from, err.Error())
					return
				}

				var coboExecuted bool
				for _, tx := range txs {
					if tx.Remark == task.TaskId {
						coboExecuted = true
						//if tx success
						if tx.Status == TxStatusSuccess {
							break
						}
						if tx.Status == TxStatusFailed {
							log.Printf("tx failed cobo id: %s", tx.CoboId)
							return
						}
					}
				}
				if !coboExecuted {
					//create tx
					_, err := m.RunTask(task)
					if err != nil {
						j, _ := json.Marshal(task)
						log.Printf("run task err: %s, task: %s", err.Error(), string(j))
						return
					}
				}
				time.Sleep(time.Second * 5)
			}
		}
	}
}

func (m *MPC) RunTask(task *Task) (string, error) {
	requestId := m.createRequestId()
	if task.operation == OperationTransfer {
		res, err := m.client.CreateTransaction(task.coin, requestId, task.amount, task.from, task.to, "", nil, task.gasPrice, task.gasLimit, task.operation, "", nil, nil, nil, task.TaskId)
		if err != nil {
			return "", wrapCoboErr(err)
		}
		coboId, _ := res.Get("cobo_id").String()
		return coboId, nil
	}
	if task.operation == OperationContractCall {
		res, err := m.client.CreateTransaction(task.coin, requestId, task.amount, task.from, task.to, "", nil, task.gasPrice, task.gasLimit, task.operation, task.callData, nil, nil, nil, task.TaskId)
		if err != nil {
			return "", wrapCoboErr(err)
		}
		coboId, _ := res.Get("cobo_id").String()
		return coboId, nil
	}
	return "", fmt.Errorf("task operation err: %d", task.operation)
}

func (t *Task) CreateTaskId() {
	t.TaskId = fmt.Sprintf("%s-%s-%s", t.to, t.amount.String(), t.callData)
	t.TaskId = strings.ToLower(t.TaskId)
}

func (m *MPC) CreateTaskTransfer(coin, from, to string, amount, gasPrice, gasLimit *big.Int) *Task {
	t := Task{
		TaskId:    "",
		coin:      coin,
		from:      from,
		to:        to,
		operation: OperationTransfer,
		callData:  "",
		amount:    amount,
		gasPrice:  gasPrice,
		gasLimit:  gasLimit,
	}
	t.CreateTaskId()
	return &t
}

func (m *MPC) CreateTaskApprove(coin, from, token, spender string, approveAmount, gasPrice, gasLimit *big.Int) *Task {
	abiReader, _ := abi.JSON(strings.NewReader(ApproveABI))
	approvePack, _ := abiReader.Pack("approve", common.HexToAddress(spender), approveAmount)

	t := Task{
		TaskId:    "",
		coin:      coin,
		from:      from,
		to:        token,
		operation: OperationContractCall,
		callData:  fmt.Sprintf("0x%s", common.Bytes2Hex(approvePack)),
		amount:    big.NewInt(0),
		gasPrice:  gasPrice,
		gasLimit:  gasLimit,
	}
	t.CreateTaskId()
	return &t
}

func (m *MPC) CreateTaskContractCall(coin, from, to, callData string, amount, gasPrice, gasLimit *big.Int) *Task {
	extraParameters := fmt.Sprintf("{\"calldata\": \"%s\"}", callData)

	t := Task{
		TaskId:    "",
		coin:      coin,
		from:      from,
		to:        to,
		operation: OperationContractCall,
		callData:  extraParameters,
		amount:    big.NewInt(0),
		gasPrice:  gasPrice,
		gasLimit:  gasLimit,
	}
	t.CreateTaskId()
	return &t
}
