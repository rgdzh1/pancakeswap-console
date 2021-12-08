package cmd

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/walletConsole/pancakeswap-console/config"
	"github.com/walletConsole/pancakeswap-console/utils"
	"log"
	"strings"
	"time"
)

type TradePairs struct {
	Refreshable
	FromSymbol string
	ToSymbol   string
}

func (tp *TradePairs) Start() {
	go tp.refreshData()
	go tp.refreshView()

}
func (tp *TradePairs) buildContent() string {
	defer func() {
		if p := recover(); p != nil {
			ShowMsg(tp.G, fmt.Sprintf("internal error: %v", p))
		}
	}()

	path, err := utils.CalculatePath(tp.FromSymbol, tp.ToSymbol, config.CF.PancakeRouter, Client)
	if err != nil {
		log.Fatal(err)
	}
	price := PairsPrice(path)
	var frombalance float64
	if strings.ToLower(tp.FromSymbol) == "bnb" {
		ethBalance := utils.EthBalance(config.CF.FromAddress, tp.Client)
		frombalance = utils.Str2Float(ethBalance)
	} else {
		erc20Token := utils.Erc20Token(config.CF.BscToken[strings.ToLower(tp.FromSymbol)], Client)
		balanceOf, err := erc20Token.BalanceOf(utils.DefaultCallOpts, common.HexToAddress(config.CF.FromAddress))
		decimals, _ := erc20Token.Decimals(utils.DefaultCallOpts)

		if err == nil {

			ethBalance := utils.ToDecimal(balanceOf.String(), int(decimals))
			frombalance = utils.Str2Float(ethBalance.String())

		}

	}
	var tobalance float64
	if strings.ToLower(tp.ToSymbol) == "bnb" {
		ethBalance := utils.EthBalance(config.CF.FromAddress, tp.Client)
		tobalance = utils.Str2Float(ethBalance)
	} else {
		erc20Token := utils.Erc20Token(config.CF.BscToken[strings.ToLower(tp.ToSymbol)], Client)
		balanceOf, err := erc20Token.BalanceOf(utils.DefaultCallOpts, common.HexToAddress(config.CF.FromAddress))
		decimals, _ := erc20Token.Decimals(utils.DefaultCallOpts)

		if err == nil {

			ethBalance := utils.ToDecimal(balanceOf.String(), int(decimals))
			tobalance = utils.Str2Float(ethBalance.String())

		}

	}

	sprintf := fmt.Sprintf("price: %f \n%s: %f \n%s: %f \n", price, tp.FromSymbol, frombalance, tp.ToSymbol, tobalance)
	return sprintf
}
func (tp *TradePairs) refreshData() {
	for {
		select {
		case <-time.After(2 * time.Second):
			sprintf := tp.buildContent()
			if sprintf == "" {
				return
			}
			tp.Content = sprintf
		}
	}
}
