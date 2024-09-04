package v2

import (
	"context"
	"crypto/tls"
	"fmt"
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2/crypto"
	"log"
	"net/http"
	"time"
)

type EvmMpcV2 struct {
	privateKey string
	client     *coboWaas2.APIClient
	walletId   string
	env        int
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

	return &mpc
}

func (m *EvmMpcV2) createRequestId() string {
	return fmt.Sprintf("cs-go-v2-%d", time.Now().UnixMilli())
}

func (m *EvmMpcV2) getCtx() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, coboWaas2.ContextEnv, m.env)
	ctx = context.WithValue(ctx, coboWaas2.ContextPortalSigner, crypto.Ed25519Signer{
		Secret: m.privateKey,
	})
	return ctx
}
