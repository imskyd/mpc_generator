package mpc_generator

const (
	OperationTransfer     = 100
	OperationContractCall = 200
)

const (
	TxStatusPendingApproval     int = 101
	TxStatusQueued              int = 201
	TxStatusPendingSignature    int = 301 //waiting a signature from tss node
	TxStatusBroadcasting        int = 401
	TxStatusBroadcastFailed     int = 402
	TxStatusPendingConfirmation int = 403
	TxStatusConfirmation        int = 501 //success after code 501
	TxStatusSuccess             int = 900
	TxStatusFailed              int = 901
)

var TextTxStatus = map[int]string{
	TxStatusPendingApproval:     "PendingApproval",
	TxStatusQueued:              "Queued",
	TxStatusPendingSignature:    "PendingSignature",
	TxStatusBroadcasting:        "Broadcasting",
	TxStatusBroadcastFailed:     "BroadcastFailed",
	TxStatusPendingConfirmation: "PendingConfirmation",
	TxStatusConfirmation:        "Confirmation",
	TxStatusSuccess:             "Success",
	TxStatusFailed:              "Failed",
}

const ApproveABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
