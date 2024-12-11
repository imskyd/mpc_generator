package v2

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2/crypto"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type EvmMpcV2 struct {
	privateKey string
	client     *coboWaas2.APIClient
	walletId   string
	env        int
	debug      bool
	logger     *logrus.Logger
	ethClient  *ethclient.Client
}

func NewEvmMpcV2(env int, privateKey, walletId, evmRpc string) *EvmMpcV2 {
	if env != coboWaas2.ProdEnv && env != coboWaas2.DevEnv {
		log.Panic("env should be coboWaas2.ProdEnv or coboWaas2.DevEnv")
	}

	mpc := EvmMpcV2{privateKey: privateKey}

	configuration := coboWaas2.NewConfiguration()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	configuration.HTTPClient = client
	apiClient := coboWaas2.NewAPIClient(configuration)

	mpc.client = apiClient
	mpc.walletId = walletId
	mpc.env = env

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	})
	mpc.logger = logger

	ethClient, _ := ethclient.Dial(evmRpc)
	mpc.ethClient = ethClient

	return &mpc
}

func (m *EvmMpcV2) createRequestId() string {
	return fmt.Sprintf("cs-go-v2-%d", time.Now().UnixMilli())
}

func (m *EvmMpcV2) SetDebug(status bool) {
	m.debug = status
}

func (m *EvmMpcV2) WalletId() string {
	return m.walletId
}

func (m *EvmMpcV2) getCtx() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, coboWaas2.ContextEnv, m.env)
	ctx = context.WithValue(ctx, coboWaas2.ContextPortalSigner, crypto.Ed25519Signer{
		Secret: m.privateKey,
	})
	return ctx
}

func (m *EvmMpcV2) printFormatLog(title string, logData interface{}) {
	jsonData, _ := json.Marshal(logData)
	m.logger.Infof("%s: %s", title, string(jsonData))
}

func (m *EvmMpcV2) formatResponse(resp *coboWaas2.CreateTransferTransaction201Response, r *http.Response, err error) (*coboWaas2.CreateTransferTransaction201Response, error) {
	if err == nil {
		return resp, nil
	}
	defer r.Body.Close()

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		return nil, fmt.Errorf("err: %s, unable to read response body: %s", err.Error(), readErr)
	}
	return nil, fmt.Errorf("err: %s, body: %s", err.Error(), string(body))
}

func (m *EvmMpcV2) formatResponseCommon(resp *http.Response, err error) (*http.Response, error) {
	if err == nil {
		return resp, nil
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("err: %s, unable to read response body: %s", err.Error(), readErr)
	}
	return nil, fmt.Errorf("err: %s, body: %s", err.Error(), string(body))
}

func (m *EvmMpcV2) EstimateGas(from, to, data, value, gasPrice, maxFee, priorityGas string) (uint64, error) {
	fromAddr := common.HexToAddress(from)
	toAddr := common.HexToAddress(to)
	if strings.HasPrefix(data, "0x") {
		data = data[2:]
	}
	dataHex := common.FromHex(data)
	gasPriceBigInt, _ := decimal.NewFromString(gasPrice)
	maxFeeBigInt, _ := decimal.NewFromString(maxFee)
	priorityGasBigInt, _ := decimal.NewFromString(priorityGas)

	valueBigInt, _ := decimal.NewFromString(value)
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(18)))
	result := valueBigInt.Mul(mul)

	msg := ethereum.CallMsg{
		From:      fromAddr,
		To:        &toAddr,
		Gas:       0,
		GasPrice:  gasPriceBigInt.BigInt(),
		GasFeeCap: maxFeeBigInt.BigInt(),
		GasTipCap: priorityGasBigInt.BigInt(),
		Value:     result.BigInt(),
		Data:      dataHex,
	}

	gas, err := m.ethClient.EstimateGas(context.Background(), msg)
	return gas, err
}
