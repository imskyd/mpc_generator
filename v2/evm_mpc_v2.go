package v2

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2/crypto"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"time"
)

type EvmMpcV2 struct {
	privateKey string
	client     *coboWaas2.APIClient
	walletId   string
	env        int
	debug      bool
	logger     *logrus.Logger
}

func NewEvmMpcV2(env int, privateKey, walletId string) *EvmMpcV2 {
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

	return &mpc
}

func (m *EvmMpcV2) createRequestId() string {
	return fmt.Sprintf("cs-go-v2-%d", time.Now().UnixMilli())
}

func (m *EvmMpcV2) SetDebug(status bool) {
	m.debug = status
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
