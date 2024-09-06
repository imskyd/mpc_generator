package example

import (
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/imskyd/mpc_generator/base"
	v2 "github.com/imskyd/mpc_generator/v2"
	"github.com/sirupsen/logrus"
	"math/big"
	"time"
)

var mpc *v2.EvmMpcV2
var logger *logrus.Logger

const env = coboWaas2.ProdEnv
const privateKey = "<PRIVATE_KEY>"
const walletId = "<WALLET_ID>"
const tokenId = "<TOKEN_ID>"

func init() {
	mpc = v2.NewEvmMpcV2(env, privateKey, walletId)

	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	})
}

func main() {
	execAddress := []string{"address1", "address2"}

	for _, eAddr := range execAddress {
		fee := v2.New1559Fee(tokenId, "300", "100", "21000")
		//approve token
		response, err := mpc.TokenApprove(tokenId, eAddr, "<TOKEN_ADDRESS>", "<SPENDER_ADDRESS>", "3000", &fee)
		if err != nil {
			logger.Errorf("Failed to approve token: %s", err.Error())
			return
		}
		if waitErr := mpc.WaitTransactionDone(response.TransactionId, 100); waitErr != nil {
			logger.Errorf("Failed to wait transaction done: %s", waitErr.Error())
			return
		}

		fee2 := v2.New1559Fee(tokenId, "300", "100", "21000")
		//contract interact
		callData, _ := base.GetAbiCallData("abi", "transfer", common.HexToAddress("address"), big.NewInt(0))
		callResp1, callErr1 := mpc.ContractCall(tokenId, eAddr, "<TOKEN_ADDRESS>", callData, "0.00001", &fee2)
		if callErr1 != nil {
			logger.Errorf("Failed to approve token: %s", callErr1.Error())
			return
		}
		if waitErr := mpc.WaitTransactionDone(callResp1.TransactionId, 100); waitErr != nil {
			logger.Errorf("Failed to wait transaction done: %s", waitErr.Error())
			return
		}
	}

}
