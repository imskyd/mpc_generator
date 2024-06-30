package mpc_generator

type Transaction struct {
	CoboId     string `json:"cobo_id"`
	RequestId  string `json:"request_id"`
	Status     int    `json:"status"`
	CoinDetail struct {
		Coin                string      `json:"coin"`
		DisplayCode         string      `json:"display_code"`
		Description         string      `json:"description"`
		Decimal             int         `json:"decimal"`
		CanDeposit          interface{} `json:"can_deposit"`
		CanWithdraw         interface{} `json:"can_withdraw"`
		ConfirmingThreshold int         `json:"confirming_threshold"`
	} `json:"coin_detail"`
	AmountDetail struct {
		Amount    string `json:"amount"`
		AbsAmount string `json:"abs_amount"`
	} `json:"amount_detail"`
	FeeDetail struct {
		FeeCoinDetail struct {
			Coin                string      `json:"coin"`
			DisplayCode         string      `json:"display_code"`
			Description         string      `json:"description"`
			Decimal             int         `json:"decimal"`
			CanDeposit          interface{} `json:"can_deposit"`
			CanWithdraw         interface{} `json:"can_withdraw"`
			ConfirmingThreshold int         `json:"confirming_threshold"`
		} `json:"fee_coin_detail"`
		GasPrice int `json:"gas_price"`
		GasLimit int `json:"gas_limit"`
		FeeUsed  int `json:"fee_used"`
	} `json:"fee_detail"`
	SourceAddresses string      `json:"source_addresses"`
	FromAddress     string      `json:"from_address"`
	ToAddress       string      `json:"to_address"`
	TxHash          string      `json:"tx_hash"`
	VoutN           int         `json:"vout_n"`
	Nonce           interface{} `json:"nonce"`
	ConfirmedNumber int         `json:"confirmed_number"`
	ReplaceCoboId   string      `json:"replace_cobo_id"`
	TransactionType int         `json:"transaction_type"`
	Operation       int         `json:"operation"`
	BlockDetail     struct {
		BlockHash   string `json:"block_hash"`
		BlockHeight int    `json:"block_height"`
		BlockTime   int    `json:"block_time"`
	} `json:"block_detail"`
	TxDetail struct {
		TxHash string `json:"tx_hash"`
	} `json:"tx_detail"`
	ExtraParameters string      `json:"extra_parameters"`
	CreatedTime     int64       `json:"created_time"`
	UpdatedTime     int64       `json:"updated_time"`
	FailedReason    interface{} `json:"failed_reason"`
	MaxPriorityFee  interface{} `json:"max_priority_fee"`
	MaxFee          interface{} `json:"max_fee"`
	ApprovalProcess struct {
		SpenderResult          int `json:"spender_result"`
		SpenderReviewThreshold int `json:"spender_review_threshold"`
		SpenderStatus          []struct {
			SpenderPerson string `json:"spender_person"`
			Status        string `json:"status"`
		} `json:"spender_status"`
		SpenderCompleteTime     int64 `json:"spender_complete_time"`
		ApproverResult          int   `json:"approver_result"`
		ApproverReviewThreshold int   `json:"approver_review_threshold"`
		ApproverStatus          []struct {
			ApprovePerson string `json:"approve_person"`
			Status        string `json:"status"`
		} `json:"approver_status"`
		ApproverCompleteTime int64 `json:"approver_complete_time"`
	} `json:"approval_process"`
	Remark            string `json:"remark"`
	Memo              string `json:"memo"`
	GasStationChildId string `json:"gas_station_child_id"`
}
