package example

import (
	"context"
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imskyd/mpc_generator/base"
	v2 "github.com/imskyd/mpc_generator/v2"
	"github.com/shopspring/decimal"
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
const rpc = "<RPC_URL>"

func init() {
	mpc = v2.NewEvmMpcV2(env, privateKey, walletId, rpc)

	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	})
}

func MulGasPrice(gasPrice *big.Int, mul float64) *big.Int {
	mulGasPriceInit, _ := decimal.NewFromString(gasPrice.String())
	mulGasPrice := mulGasPriceInit.Mul(decimal.NewFromFloat(mul))
	return mulGasPrice.BigInt()
}

func main() {
	execAddress := []string{"address1", "address2"}
	client, _ := ethclient.Dial("<RPC_URL>")

	for _, eAddr := range execAddress {
		gasPrice, err := client.SuggestGasPrice(context.Background())
		maxFee := MulGasPrice(gasPrice, 2)
		priority := MulGasPrice(gasPrice, 0.1)
		fee := v2.New1559Fee(tokenId, maxFee.String(), priority.String(), "27000")
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

		gasPrice2, _ := client.SuggestGasPrice(context.Background())
		maxFee2 := MulGasPrice(gasPrice2, 2)
		priority2 := MulGasPrice(gasPrice2, 0.1)
		//contract interact
		callData, _ := base.GetAbiCallData("abi", "transfer", common.HexToAddress("address"), big.NewInt(0))
		gasLimit2, err := mpc.EstimateGas(eAddr, "<TOKEN_ADDRESS>", callData, "0.00001", "", maxFee2.String(), priority2.String())
		if err != nil {
			logger.Errorf("Failed to approve token: %s", err.Error())
			return
		}
		fee2 := v2.New1559Fee(tokenId, maxFee2.String(), priority2.String(), decimal.NewFromUint64(gasLimit2).String())
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
