/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"
	"github.com/walletConsole/pancakeswap-console/config"
	"github.com/walletConsole/pancakeswap-console/utils"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"
)

var (
	client  *ethclient.Client
	account Account
	preView *gocui.View
)

type Refreshable struct {
	client   *ethclient.Client
	Rapid    float64
	Fast     float64
	Standard float64
	g        *gocui.Gui
	view     string
	content  chan string
}

func (tp *Refreshable) refreshView() {
	for {
		sprintf := <-tp.content
		tp.g.Update(func(g *gocui.Gui) error {
			v, err := g.View(tp.view)
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintln(v, sprintf)
			return nil
		})
	}
}

type Account struct {
	Refreshable
	address string
}
type TradePairs struct {
	Refreshable
	fromSymbol string
	toSymbol   string
}

func (tp *Account) start() {
	go tp.refreshData()
	go tp.refreshView()

}
func (tp *Account) buildContent() string {
	defer func() {
		if p := recover(); p != nil {
			showMsg(tp.g, fmt.Sprintf("internal error: %v", p))
		}
	}()
	ethBalance := utils.EthBalance(tp.address, tp.client)
	balance := utils.Str2Float(ethBalance)
	var path []common.Address
	path = append(path, common.HexToAddress(config.CF.BscToken[strings.ToLower("WBNB")]))
	path = append(path, common.HexToAddress(config.CF.BscToken[strings.ToLower("USDT")]))

	price := pairsPrice(path)

	sprintf := fmt.Sprintf("%s: %f (BNB) (%f)", tp.address, balance, price)

	erc20Token := utils.Erc20Token(config.CF.BscToken[strings.ToLower("usdt")], tp.client)
	balanceOf, err := erc20Token.BalanceOf(utils.DefaultCallOpts, common.HexToAddress(config.CF.FromAddress))
	decimals, _ := erc20Token.Decimals(utils.DefaultCallOpts)
	usdtbalance := 0.0
	if err == nil {
		if balanceOf.Uint64() > 0 {
			ethBalance := utils.ToDecimal(balanceOf.String(), int(decimals))
			usdtbalance = utils.Str2Float(ethBalance.String())

		}
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	tp.Refreshable.Standard = utils.Str2Float(utils.ToDecimal(gasPrice, 9).String())
	tp.Refreshable.Rapid = tp.Refreshable.Standard * 1.5
	tp.Refreshable.Fast = tp.Refreshable.Standard * 1.2
	sprintf = sprintf + fmt.Sprintf("\t%f (USDT)", usdtbalance)
	sprintf = sprintf + fmt.Sprintf("\t\t\tGas: %f(gwei)\t%f(gwei)\t%f(gwei) ", tp.Refreshable.Rapid, tp.Refreshable.Fast, tp.Refreshable.Standard)
	return sprintf
}
func (tp *Account) refreshData() {
	for {
		select {
		case <-time.After(2 * time.Second):
			sprintf := tp.buildContent()
			if sprintf == "" {
				return
			}
			tp.content <- sprintf
		}
	}
}
func (tp *TradePairs) start() {
	go tp.refreshData()
	go tp.refreshView()

}
func (tp *TradePairs) buildContent() string {
	defer func() {
		if p := recover(); p != nil {
			showMsg(tp.g, fmt.Sprintf("internal error: %v", p))
		}
	}()

	path, err := utils.CalculatePath(tp.fromSymbol, tp.toSymbol, config.CF.PancakeRouter, client)
	if err != nil {
		log.Fatal(err)
	}
	price := pairsPrice(path)
	var frombalance float64
	if strings.ToLower(tp.fromSymbol) == "bnb" {
		ethBalance := utils.EthBalance(config.CF.FromAddress, tp.client)
		frombalance = utils.Str2Float(ethBalance)
	} else {
		erc20Token := utils.Erc20Token(config.CF.BscToken[strings.ToLower(tp.fromSymbol)], client)
		balanceOf, err := erc20Token.BalanceOf(utils.DefaultCallOpts, common.HexToAddress(config.CF.FromAddress))
		decimals, _ := erc20Token.Decimals(utils.DefaultCallOpts)

		if err == nil {
			if balanceOf.Uint64() > 0 {
				ethBalance := utils.ToDecimal(balanceOf.String(), int(decimals))
				frombalance = utils.Str2Float(ethBalance.String())
			}
		}

	}
	var tobalance float64
	if strings.ToLower(tp.toSymbol) == "bnb" {
		ethBalance := utils.EthBalance(config.CF.FromAddress, tp.client)
		tobalance = utils.Str2Float(ethBalance)
	} else {
		erc20Token := utils.Erc20Token(config.CF.BscToken[strings.ToLower(tp.toSymbol)], client)
		balanceOf, err := erc20Token.BalanceOf(utils.DefaultCallOpts, common.HexToAddress(config.CF.FromAddress))
		decimals, _ := erc20Token.Decimals(utils.DefaultCallOpts)

		if err == nil {
			if balanceOf.Uint64() > 0 {
				ethBalance := utils.ToDecimal(balanceOf.String(), int(decimals))
				tobalance = utils.Str2Float(ethBalance.String())
			}
		}

	}

	sprintf := fmt.Sprintf("price: %f \n%s: %f \n%s: %f \n", price, tp.fromSymbol, frombalance, tp.toSymbol, tobalance)
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
			tp.content <- sprintf
		}
	}
}

var monitorSwapCmd = &cobra.Command{
	Use:   "swap",
	Short: "run pancake swap",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client = utils.EstimateClient("bsc")
		g, err := gocui.NewGui(gocui.OutputNormal)
		if err != nil {
			log.Panicln(err)
		}
		defer g.Close()
		g.Cursor = true
		g.Mouse = true
		g.Highlight = true
		//g.Cursor = true
		g.SelFgColor = gocui.ColorGreen

		g.SetManagerFunc(layout)

		if err := keybindings(g); err != nil {
			log.Panicln(err)
		}
		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}

		//for key, value := range config.CF.SyrupPool {
		//	monitorSyrup(key,value)
		//}
	},
}

func GetAmountsOut(fromSymbol, toSymbol, amount string) float64 {
	var swapFromTokenAddress string
	if strings.ToLower(fromSymbol) == "bnb" {
		swapFromTokenAddress = config.BscToken("wbnb")
	} else {
		swapFromTokenAddress = config.BscToken(fromSymbol)
	}
	var swapToTokenAddress string
	if strings.ToLower(toSymbol) == "bnb" {
		swapToTokenAddress = config.BscToken("wbnb")
	} else {
		swapToTokenAddress = config.BscToken(fromSymbol)
	}

	path, err := utils.CalculatePath(fromSymbol, toSymbol, config.CF.PancakeRouter, client)
	if err != nil {
		log.Fatal(err)
	}

	value := utils.Str2Float(amount)
	swapFrom := utils.Erc20Token(swapFromTokenAddress, client)
	swapFromDecimals, _ := swapFrom.Decimals(utils.DefaultCallOpts)
	amountInBig := utils.ToWei(value, int(swapFromDecimals))

	swapTo := utils.Erc20Token(swapToTokenAddress, client)
	swapToDecimals, _ := swapTo.Decimals(utils.DefaultCallOpts)

	amountsOut := utils.GetAmountsOut(config.CF.PancakeRouter, path, amountInBig, client)
	//log.Printf("%v", amountsOut)
	tokenAmountOutFloat := SlippageAmount(amountsOut, swapToDecimals, err)
	return tokenAmountOutFloat
}
func SlippageAmount(amountsOut *big.Int, swapToDecimals uint8, err error) float64 {
	tokenAmountOutStr := utils.ToDecimal(amountsOut.String(), int(swapToDecimals)).String()

	tokenAmountOutFloat, err := strconv.ParseFloat(tokenAmountOutStr, 64)
	if err != nil {
		panic(err)
	}
	tokenAmountOutFloat = tokenAmountOutFloat - tokenAmountOutFloat*0.01
	return tokenAmountOutFloat
}

func Swap(fromSymbol, toSymbol, amount string) string {
	var swapFromTokenAddress string
	if strings.ToLower(fromSymbol) == "bnb" {
		swapFromTokenAddress = config.BscToken("wbnb")
	} else {
		swapFromTokenAddress = config.BscToken(fromSymbol)
	}
	var swapToTokenAddress string
	if strings.ToLower(toSymbol) == "bnb" {
		swapToTokenAddress = config.BscToken("wbnb")
	} else {
		swapToTokenAddress = config.BscToken(fromSymbol)
	}

	path, err := utils.CalculatePath(fromSymbol, toSymbol, config.CF.PancakeRouter, client)
	if err != nil {
		log.Fatal(err)
	}

	value := utils.Str2Float(amount)
	swapFrom := utils.Erc20Token(swapFromTokenAddress, client)
	swapFromDecimals, _ := swapFrom.Decimals(utils.DefaultCallOpts)
	amountInBig := utils.ToWei(value, int(swapFromDecimals))

	swapTo := utils.Erc20Token(swapToTokenAddress, client)
	swapToDecimals, _ := swapTo.Decimals(utils.DefaultCallOpts)

	amountsOut := utils.GetAmountsOut(config.CF.PancakeRouter, path, amountInBig, client)
	//log.Printf("%v", amountsOut)
	tokenAmountOutFloat := SlippageAmount(amountsOut, swapToDecimals, err)
	uniswapv2router02 := utils.GetUni2router2(config.CF.PancakeRouter, client)

	fromAddress := common.HexToAddress(config.CF.FromAddress)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	suggestGasPrice, err := client.SuggestGasPrice(context.Background())
	gasPrice := big.NewInt(int64(1)).Add(suggestGasPrice, big.NewInt(int64(1)))
	privateKey, err := crypto.HexToECDSA(config.CF.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	opts := bind.NewKeyedTransactor(privateKey)
	opts.Nonce = big.NewInt(int64(nonce))
	opts.GasLimit = uint64(255164) // in units
	opts.GasPrice = gasPrice

	deadline := &big.Int{}
	deadline.SetInt64(time.Now().Add(10 * time.Minute).Unix())

	var transaction *types.Transaction
	if strings.ToLower(fromSymbol) == "bnb" {

		opts.Value = amountInBig // in wei

		transaction, err = uniswapv2router02.SwapExactETHForTokens(opts, utils.ToWei(tokenAmountOutFloat, int(swapToDecimals)), path, fromAddress, deadline)
		if err != nil {
			log.Fatal(err)
		}
	} else if strings.ToLower(toSymbol) == "bnb" {

		allowance, err := swapFrom.Allowance(utils.DefaultCallOpts, fromAddress, common.HexToAddress(config.CF.PancakeRouter))
		if allowance.Uint64() < amountInBig.Uint64() {

			utils.ApproveErc20AndAwait(config.CF.PrivateKey, config.CF.PancakeRouter, config.BscToken(fromSymbol), -1, uint64(utils.Str2Int(utils.ToDecimal(gasPrice, 9).String())), client)

			opts.Nonce = big.NewInt(int64(nonce + 1))
		}

		transaction, err = uniswapv2router02.SwapExactTokensForETH(opts, amountInBig, utils.ToWei(tokenAmountOutFloat, int(swapToDecimals)), path, fromAddress, deadline)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		allowance, err := swapFrom.Allowance(utils.DefaultCallOpts, fromAddress, common.HexToAddress(config.CF.PancakeRouter))
		if allowance.Uint64() < amountInBig.Uint64() {
			utils.ApproveErc20AndAwait(config.CF.PrivateKey, config.CF.PancakeRouter, config.BscToken(fromSymbol), -1, uint64(utils.Str2Int(utils.ToDecimal(gasPrice, 9).String())), client)
			opts.Nonce = big.NewInt(int64(nonce + 1))
		}

		transaction, err = uniswapv2router02.SwapExactTokensForTokens(opts, amountInBig, utils.ToWei(tokenAmountOutFloat, int(swapToDecimals)), path, fromAddress, deadline)
		if err != nil {
			log.Fatal(err)
		}
	}

	return transaction.Hash().String()
}

func pairsPrice(path []common.Address) float64 {
	//var path []common.Address
	tokenAddress := path[0]
	//client := utils.EstimateClient(Chain)
	erc20Token := utils.Erc20Token(tokenAddress.Hex(), client)
	decimals, _ := erc20Token.Decimals(utils.DefaultCallOpts)
	newInt := big.NewInt(int64(utils.EthUnit))
	if decimals == 6 {
		newInt = big.NewInt(int64(utils.WeiUnit6))
	} else if decimals == 8 {
		newInt = big.NewInt(int64(utils.WeiUnit8))
	}

	price := utils.PriceFromSwap(config.CF.PancakeRouter, path, newInt, client)
	return price
}

func init() {
	rootCmd.AddCommand(monitorSwapCmd)
}
