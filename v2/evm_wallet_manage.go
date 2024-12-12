package v2

import (
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

func (m *EvmMpcV2) WalletCreateAddress(walletId, chainId string, walletNum int32) ([]coboWaas2.AddressInfo, error) {
	createAddressRequest := *coboWaas2.NewCreateAddressRequest(chainId, walletNum)
	addresses, resp, err := m.client.WalletsAPI.CreateAddress(m.getCtx(), walletId).CreateAddressRequest(createAddressRequest).Execute()
	_, errResp := m.formatResponseCommon(resp, err)
	return addresses, errResp
}

func (m *EvmMpcV2) ListTokenBalancesForAddress(walletId, address string) (*coboWaas2.ListTokenBalancesForAddress200Response, error) {
	apiResp, resp, err := m.client.WalletsAPI.ListTokenBalancesForAddress(m.getCtx(), walletId, address).Execute()
	_, errResp := m.formatResponseCommon(resp, err)
	return apiResp, errResp
}
