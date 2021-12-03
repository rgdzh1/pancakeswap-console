# Pancakeswap Console
A Console Application Use Pancakeswap To Swap Token

## Why
It is a pity that some areas do not have access to the functions of PancakeSwap.
I developed this tool to help these $People$.

I believe $People$ have the right to access the blockchain.

## Building on Linux

Dependencies:

   * go >= 1.9
   
Clone & compile:
    
    git clone https://github.com/walletConsole/pancakeswap-console.git
    go mod tidy
    go build

## Features
* RealTime Prices      [x]
* Simple Swap          [x]
* Syrup Pool           [ ]
* IFO                  [ ]
* Batch Transfer       [ ]
* Batch Deposit Syrup Pool [ ]

## Use
rename tool-conf.yaml.sample to tool-conf.yaml. Add Address And Private Key.

Your private key is secure. You can check the code.

```yaml
Bsc:
  RpcUrl: "https://bsc-dataseed1.defibit.io/"

# u can check code , u private key is safe
FromAddress: "xxx"
PrivateKey:  "xxxx"

#  Token to Swap,Do not add too many Token
PriceToken:
  - BNB-USDT
  - CAKE-USDT

# Bsc Token Info
BscToken: {
  CAKE: "0x0e09fabb73bd3ade0a17ecc321fd13a19e81ce82",
  DAR: "0x23CE9e926048273eF83be0A3A8Ba9Cb6D45cd978",
  WBNB: "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c",
  BUSD: "0xe9e7cea3dedca5984780bafc599bd69add087d56",
  USDT: "0x55d398326f99059ff775485246999027b3197955",
  QI: "0x8729438eb15e2c8b576fcc6aecda6a148776c0f5",
  KART: "0x8bdd8dbcbdf0c066ca5f3286d33673aa7a553c10",
  PORTO: "0x49f2145d6366099e13b10fbf80646c0f377ee7f6",
  ETERNAL: "0xD44FD09d74cd13838F137B590497595d6b3FEeA4",
  SFUND: "0x477bc8d23c634c154061869478bce96be6045d12",
  ZOO: "0x1D229B958D5DDFca92146585a8711aECbE56F095",
  QUIDD: "0x7961ade0a767c0e5b67dd1a1f78ba44f727642ed",
  SANTOS: "0xa64455a4553c9034236734faddaddbb64ace4cc7",
}

```
run

     ./pancakeswap-console swap

Update Price Every Second
![image](https://github.com/walletConsole/pancakeswap-console/blob/master/1.jpg)

### Swap 
For Swap BNB <> USDT
Mouse Left BNB Sell to USDT ,Mouse Left USDT Buy to BNB

![image](https://github.com/walletConsole/pancakeswap-console/blob/master/2.jpg)

Input Swap Amount

![image](https://github.com/walletConsole/pancakeswap-console/blob/master/3.jpg)

Confirm Swap

![image](https://github.com/walletConsole/pancakeswap-console/blob/master/4.jpg)

Sending

![image](https://github.com/walletConsole/pancakeswap-console/blob/master/5.jpg)



### Note
* The price may be a little different from that of the official website.Because using route path  <token> -> WBNB -> USDT
* Don't set up too many trading pairs, there will be problems with rpc node.


### Donations

BNB/CAKE/USDT: 0xec5fa25e37dfa8fa42210a94cbc8a61c7fd3751c

email: irmakgu40@gmail.com