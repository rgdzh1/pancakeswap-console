package cmd

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jroimartin/gocui"
	"github.com/walletConsole/pancakeswap-console/config"
	"github.com/walletConsole/pancakeswap-console/utils"
	"math/big"
	"strings"
	"time"
)

type Refreshable struct {
	Client *ethclient.Client
	//Rapid    float64
	//Fast     float64
	//Standard float64
	G       *gocui.Gui
	View    string
	Content string
}

func (tp *Refreshable) refreshView() {
	for {
		select {
		case <-time.After(1 * time.Second):
			//sprintf := <-tp.Content
			tp.G.Update(func(g *gocui.Gui) error {
				v, err := g.View(tp.View)
				if err != nil {
					return err
				}
				v.Clear()
				fmt.Fprintln(v, tp.Content)
				return nil
			})
		}
	}

}

type Account struct {
	Refreshable
	Address     string
	Balance     float64
	BnbPrice    float64
	UsdtBalance float64
}

func (tp *Account) Start() {
	go tp.refreshData()
	go tp.refreshView()

}
func (tp *Account) output() string {
	defer func() {
		if p := recover(); p != nil {
			ShowMsg(tp.G, fmt.Sprintf("internal error: %v", p))
		}
	}()
	sprintf := fmt.Sprintf("%s: %f (BNB)\t%f (USDT)", tp.Address, tp.Balance, tp.BnbPrice*tp.Balance)
	sprintf = sprintf + fmt.Sprintf("\tUSDT:%f", tp.UsdtBalance)
	//sprintf = sprintf + fmt.Sprintf("\t\t\tGas: %f(gwei)\t%f(gwei)\t%f(gwei) ", tp.Refreshable.Rapid, tp.Refreshable.Fast, tp.Refreshable.Standard)
	return sprintf
}
func (tp *Account) refreshData() {
	defer func() {
		if p := recover(); p != nil {
			ShowMsg(tp.G, fmt.Sprintf("internal error: %v", p))
		}
	}()

	for {
		select {
		case <-time.After(2 * time.Second):

			ethBalance := utils.EthBalance(tp.Address, tp.Client)
			balance := utils.Str2Float(ethBalance)
			tp.Balance = balance

			var path []common.Address
			path = append(path, common.HexToAddress(config.CF.BscToken[strings.ToLower("WBNB")]))
			path = append(path, common.HexToAddress(config.CF.BscToken[strings.ToLower("USDT")]))

			tp.BnbPrice = PairsPrice(path)
			erc20Token := utils.Erc20Token(config.CF.BscToken[strings.ToLower("usdt")], tp.Client)
			balanceOf, err := erc20Token.BalanceOf(utils.DefaultCallOpts, common.HexToAddress(config.CF.FromAddress))

			decimals, _ := erc20Token.Decimals(utils.DefaultCallOpts)
			tp.UsdtBalance = 0.0
			if err == nil {
				ethBalance := utils.ToDecimal(balanceOf.String(), int(decimals))
				tp.UsdtBalance = utils.Str2Float(ethBalance.String())
			}
			sprintf := tp.output()
			if sprintf == "" {
				return
			}
			tp.Content = sprintf
		}
	}
}
func PairsPrice(path []common.Address) float64 {
	//var path []common.Address
	tokenAddress := path[0]
	//client := utils.EstimateClient(Chain)
	erc20Token := utils.Erc20Token(tokenAddress.Hex(), Client)
	decimals, _ := erc20Token.Decimals(utils.DefaultCallOpts)
	newInt := big.NewInt(int64(utils.EthUnit))
	if decimals == 6 {
		newInt = big.NewInt(int64(utils.WeiUnit6))
	} else if decimals == 8 {
		newInt = big.NewInt(int64(utils.WeiUnit8))
	}

	price := utils.PriceFromSwap(config.CF.PancakeRouter, path, newInt, Client)
	return price
}
