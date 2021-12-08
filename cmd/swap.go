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
	"github.com/walletConsole/pancakeswap-console/config"
	"github.com/walletConsole/pancakeswap-console/utils"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"
)

var (
	Client  *ethclient.Client
	account Account
	preView *gocui.View
)

//var monitorSwapCmd = &cobra.Command{
//	Use:   "swap",
//	Short: "run pancake swap",
//	Long:  ``,
//	Run: func(cmd *cobra.Command, args []string) {
//		Client = utils.EstimateClient("bsc")
//		g, err := gocui.NewGui(gocui.OutputNormal)
//		if err != nil {
//			log.Panicln(err)
//		}
//		defer g.Close()
//		g.Cursor = true
//		g.Mouse = true
//		g.Highlight = true
//		//g.Cursor = true
//		g.SelFgColor = gocui.ColorGreen
//
//		g.SetManagerFunc(layout)
//
//		if err := keybindings(g); err != nil {
//			log.Panicln(err)
//		}
//		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
//			log.Panicln(err)
//		}
//
//		//for key, value := range config.CF.SyrupPool {
//		//	monitorSyrup(key,value)
//		//}
//	},
//}
//
func GetAmountsOut(fromSymbol, toSymbol, amount string) (float64, string) {
	var swapFromTokenAddress string
	if strings.ToLower(fromSymbol) == "bnb" {
		swapFromTokenAddress = config.BscToken("wbnb")
	} else {
		swapFromTokenAddress = config.BscToken(strings.ToLower(fromSymbol))
	}

	var swapToTokenAddress string
	if strings.ToLower(toSymbol) == "bnb" {
		swapToTokenAddress = config.BscToken("wbnb")
	} else {
		swapToTokenAddress = config.BscToken(strings.ToLower(toSymbol))
	}

	path, err := utils.CalculatePath(fromSymbol, toSymbol, config.CF.PancakeRouter, Client)
	if err != nil {
		log.Fatal(err)
	}
	var pathStr string
	for i, pa := range path {
		paToken := utils.Erc20Token(pa.String(), Client)
		symbol, _ := paToken.Symbol(utils.DefaultCallOpts)
		if i < len(path)-1 {
			pathStr = pathStr + fmt.Sprintf("%s -> ", symbol)
		} else {
			pathStr = pathStr + fmt.Sprintf("%s", symbol)
		}

	}

	value := utils.Str2Float(amount)
	swapFrom := utils.Erc20Token(swapFromTokenAddress, Client)
	swapFromDecimals, _ := swapFrom.Decimals(utils.DefaultCallOpts)
	amountInBig := utils.ToWei(value, int(swapFromDecimals))

	swapTo := utils.Erc20Token(swapToTokenAddress, Client)
	swapToDecimals, _ := swapTo.Decimals(utils.DefaultCallOpts)

	amountsOut := utils.GetAmountsOut(config.CF.PancakeRouter, path, amountInBig, Client)

	tokenAmountOutFloat := SlippageAmount(amountsOut, swapToDecimals, err)
	return tokenAmountOutFloat, pathStr
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

	path, err := utils.CalculatePath(fromSymbol, toSymbol, config.CF.PancakeRouter, Client)
	if err != nil {
		log.Fatal(err)
	}

	value := utils.Str2Float(amount)
	swapFrom := utils.Erc20Token(swapFromTokenAddress, Client)
	swapFromDecimals, _ := swapFrom.Decimals(utils.DefaultCallOpts)
	amountInBig := utils.ToWei(value, int(swapFromDecimals))

	swapTo := utils.Erc20Token(swapToTokenAddress, Client)
	swapToDecimals, _ := swapTo.Decimals(utils.DefaultCallOpts)

	amountsOut := utils.GetAmountsOut(config.CF.PancakeRouter, path, amountInBig, Client)
	//log.Printf("%v", amountsOut)
	tokenAmountOutFloat := SlippageAmount(amountsOut, swapToDecimals, err)
	uniswapv2router02 := utils.GetUni2router2(config.CF.PancakeRouter, Client)

	fromAddress := common.HexToAddress(config.CF.FromAddress)
	nonce, err := Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	suggestGasPrice, err := Client.SuggestGasPrice(context.Background())
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

			utils.ApproveErc20AndAwait(config.CF.PrivateKey, config.CF.PancakeRouter, config.BscToken(fromSymbol), -1, uint64(utils.Str2Int(utils.ToDecimal(gasPrice, 9).String())), Client)

			opts.Nonce = big.NewInt(int64(nonce + 1))
		}

		transaction, err = uniswapv2router02.SwapExactTokensForETH(opts, amountInBig, utils.ToWei(tokenAmountOutFloat, int(swapToDecimals)), path, fromAddress, deadline)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		allowance, err := swapFrom.Allowance(utils.DefaultCallOpts, fromAddress, common.HexToAddress(config.CF.PancakeRouter))
		if allowance.Uint64() < amountInBig.Uint64() {
			utils.ApproveErc20AndAwait(config.CF.PrivateKey, config.CF.PancakeRouter, config.BscToken(fromSymbol), -1, uint64(utils.Str2Int(utils.ToDecimal(gasPrice, 9).String())), Client)
			opts.Nonce = big.NewInt(int64(nonce + 1))
		}

		transaction, err = uniswapv2router02.SwapExactTokensForTokens(opts, amountInBig, utils.ToWei(tokenAmountOutFloat, int(swapToDecimals)), path, fromAddress, deadline)
		if err != nil {
			log.Fatal(err)
		}
	}

	return transaction.Hash().String()
}

func init() {
	//rootCmd.AddCommand(monitorSwapCmd)
}
