package cmd

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/walletConsole/pancakeswap-console/abi/cakevault"
	"github.com/walletConsole/pancakeswap-console/abi/pancake"
	"github.com/walletConsole/pancakeswap-console/config"
	"github.com/walletConsole/pancakeswap-console/utils"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"math/big"
	"strings"
	"time"
)

type SyrupPool struct {
	Refreshable
	PoolName    string
	PoolAddress string
	Total       float64
	Reward      float64
	RewardUsdt  float64
	Stake       float64
}

func (tp *SyrupPool) Start() {
	go tp.refreshData()
	go tp.refreshView()

}
func (tp *SyrupPool) refreshData() {
	defer func() {
		if p := recover(); p != nil {
			ShowMsg(tp.G, fmt.Sprintf("internal error: %v", p))
		}
	}()
	for {
		select {
		case <-time.After(10 * time.Second):
			if strings.ToLower(tp.PoolName) == "cake-auto" {
				CakeAutoSyrupPool(tp)
			} else {
				otherSyrupPool(tp)
			}

			sprintf := tp.output()
			if sprintf == "" {
				return
			}
			tp.Content = sprintf
		}
	}
}

func CakeAutoSyrupPool(tp *SyrupPool) {
	erc20Token := utils.Erc20Token(config.CF.BscToken[strings.ToLower("cake")], Client)
	decimals, _ := erc20Token.Decimals(utils.DefaultCallOpts)
	tp.PoolAddress = config.CF.SyrupPool[strings.ToLower("cake-auto")]

	toAddress := common.HexToAddress(tp.PoolAddress)
	pancakeStakePool, err := cakevault.NewCakevault(toAddress, Client)
	if err != nil {
		log.Fatal(err)
	}
	path, err := utils.CalculatePath("cake", "USDT", config.CF.PancakeRouter, Client)
	if err != nil {
		log.Fatal(err)
	}
	price := PairsPrice(path)

	fromAddress := common.HexToAddress(config.CF.FromAddress)

	userInfo, err := pancakeStakePool.UserInfo(utils.DefaultCallOpts, fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	totalAmountBig, err := pancakeStakePool.BalanceOf(utils.DefaultCallOpts)

	totalShares, err := pancakeStakePool.TotalShares(utils.DefaultCallOpts)
	shares := userInfo.Shares

	mul := big.NewInt(0).Mul(totalAmountBig, shares)
	nowAmountBig := big.NewInt(0).Div(mul, totalShares)

	stakeAmountStr := utils.ToDecimal(userInfo.CakeAtLastUserAction.String(), int(decimals)).String()
	stakeFloat := utils.Str2Float(stakeAmountStr)

	nowAmountStr := utils.ToDecimal(nowAmountBig.String(), int(decimals)).String()
	nowAmount := utils.Str2Float(nowAmountStr)

	//rewardAmountStr, stakeAmountStr := SyrupPoolCheckReward(tp.PoolAddress, pancakeStakePool, decimals)
	//cakeErc20Token := utils.Erc20Token(config.CF.BscToken["cake"], Client)
	//totalAmountBig, err := cakeErc20Token.BalanceOf(utils.DefaultCallOpts, common.HexToAddress(tp.PoolAddress))
	totalAmount := utils.Str2Float(utils.ToDecimal(totalAmountBig.String(), 18).String())
	rewardFloat := nowAmount - stakeFloat
	//ethBalance := utils.Str2Float(ethBalanceStr)

	tp.Total = totalAmount
	tp.Reward = rewardFloat
	tp.Stake = stakeFloat
	tp.RewardUsdt = rewardFloat * price
}
func otherSyrupPool(tp *SyrupPool) {
	erc20Token := utils.Erc20Token(config.CF.BscToken[strings.ToLower(tp.PoolName)], Client)
	decimals, _ := erc20Token.Decimals(utils.DefaultCallOpts)
	tp.PoolAddress = config.CF.SyrupPool[strings.ToLower(tp.PoolName)]

	toAddress := common.HexToAddress(tp.PoolAddress)
	pancakeStakePool, err := pancake.NewPancake(toAddress, Client)
	if err != nil {
		log.Fatal(err)
	}
	path, err := utils.CalculatePath(tp.PoolName, "USDT", config.CF.PancakeRouter, Client)
	if err != nil {
		log.Fatal(err)
	}
	price := PairsPrice(path)

	//rewardPerBlock, err := pancakeStakePool.RewardPerBlock(utils.DefaultCallOpts)
	//rewardPerBlockFloat := utils.Str2Float(utils.ToDecimal(rewardPerBlock.String(), int(decimals)).String())
	//blockRewardUsdt := price * rewardPerBlockFloat
	rewardAmountStr, stakeAmountStr := SyrupPoolCheckReward(config.CF.FromAddress, pancakeStakePool, decimals)
	cakeErc20Token := utils.Erc20Token(config.CF.BscToken["cake"], Client)
	totalAmountBig, err := cakeErc20Token.BalanceOf(utils.DefaultCallOpts, common.HexToAddress(tp.PoolAddress))

	totalAmount := utils.Str2Float(utils.ToDecimal(totalAmountBig.String(), 18).String())
	rewardFloat := utils.Str2Float(rewardAmountStr)
	stakeFloat := utils.Str2Float(stakeAmountStr)
	//ethBalance := utils.Str2Float(ethBalanceStr)

	tp.Total = totalAmount
	tp.Reward = rewardFloat
	tp.Stake = stakeFloat
	tp.RewardUsdt = rewardFloat * price
}

func (tp *SyrupPool) output() string {
	//tokenAddress := common.HexToAddress(config.CF.BscToken[strings.ToLower(tp.PoolName)])
	p := message.NewPrinter(language.English)
	var sprintf string
	if strings.ToLower(tp.PoolName) == "cake-auto" {
		CakeAutoSyrupPool(tp)
		sprintf = fmt.Sprintf("Total: %s CAKE \nYour Stake: %s \nEarned: %s (CAKE)\nEarned: %s (USDT)",
			p.Sprintf("%.2f", tp.Total),
			p.Sprintf("%.2f", tp.Stake),
			p.Sprintf("%.2f", tp.Reward),
			p.Sprintf("%.2f", tp.RewardUsdt))

	} else {
		sprintf = fmt.Sprintf("Total: %s cake \nYour Stake: %s \nEarned: %s (%s)\nEarned: %s (USDT)",
			p.Sprintf("%.2f", tp.Total),
			p.Sprintf("%.2f", tp.Stake),
			p.Sprintf("%.2f", tp.Reward),
			tp.PoolName,
			p.Sprintf("%.2f", tp.RewardUsdt))

	}
	return sprintf
}

func SyrupPoolCheckReward(fromAddresss string, pancakeStakePool *pancake.Pancake, decimals uint8) (string, string) {
	fromAddress := common.HexToAddress(fromAddresss)

	pendingReward, err := pancakeStakePool.PendingReward(utils.DefaultCallOpts, fromAddress)

	if err != nil {
		log.Fatal(err)
	}
	userInfo, err := pancakeStakePool.UserInfo(utils.DefaultCallOpts, fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	return utils.ToDecimal(pendingReward.String(), int(decimals)).String(), utils.ToDecimal(userInfo.Amount.String(), 18).String()

}
