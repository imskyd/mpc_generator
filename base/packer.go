package base

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"strings"
)

const ApproveABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

func GetApproveCallData(spender, approveAmount string) (string, error) {
	abiReader, _ := abi.JSON(strings.NewReader(ApproveABI))
	approvePack, packErr := abiReader.Pack("approve", common.HexToAddress(spender), approveAmount)
	if packErr != nil {
		return "", packErr
	}
	return fmt.Sprintf("0x%s", common.Bytes2Hex(approvePack)), nil
}

func GetAbiCallData(abiJson, name string, args ...interface{}) (string, error) {
	abiReader, _ := abi.JSON(strings.NewReader(abiJson))
	approvePack, packErr := abiReader.Pack(name, args...)
	if packErr != nil {
		return "", packErr
	}
	return fmt.Sprintf("0x%s", common.Bytes2Hex(approvePack)), nil
}
