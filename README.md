# Pancakeswap Console
A Pancakeswap Application 

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
* Prices      [x]
* Syrup Pool           [x]
* IFO                  [ ]
* Batch Transfer       [ ]
* Batch Deposit Syrup Pool [ ]





## Use
rename tool-conf.yaml.sample to tool-conf.yaml. Set You Address


```yaml
Bsc:
  RpcUrl: "https://bsc-dataseed1.defibit.io/"
  WsUrl:  "wss://bsc-ws-node.nariox.org:443"
  
# set u adddress to Analysis
FromAddress: "XXXX"

PancakeRouter: "0x10ed43c718714eb63d5aa57b78b54704e256024e"

PriceToken:
  - BNB-USDT
  - CAKE-USDT
  - DAR-USDT
  - XCV-USDT
  - SANTOS-USDT

ShowSyrupPool:
  - CAKE-AUTO
  - SANTOS
  - QUIDD
  - XCV

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
  XCV: "0x4be63a9b26EE89b9a3a13fd0aA1D0b2427C135f8",
}
SyrupPool: {
  CAKE-AUTO: "0xa80240Eb5d7E05d3F250cF000eEc0891d00b51CC",
  CAKE: "0xbb472601b3cb32723d0755094ba80b73f67f2af3",
  AXS: "0xbb472601b3cb32723d0755094ba80b73f67f2af3",
  QI: "0xbd52ef04DB1ad1c68A8FA24Fa71f2188978ba617",
  KART: "0x73bB10B89091f15e8FeD4d6e9EBa6415df6acb21",
  PORTO: "0xdD52FAB121376432DBCBb47592742F9d86CF8952",
  ETERNAL: "0xc28c400F2B675b25894FA632205ddec71E432288",
  SFUND: "0x7f103689cabe17c2f70da6faa298045d72a943b8",
  DAR: "0x9b861A078B2583373A7a3EEf815bE1A39125Ae08",
  ZOO: "0x2EfE8772EB97B74be742d578A654AB6C95bF18db",
  QUIDD: "0xd97ee2bfe79a4d4ab388553411c462fbb536a88c",
  SANTOS: "0x0914b2d9d4dd7043893def53ecfc0f1179f87d5c",
  XCV: "0xf1fa41f593547e406a203b681df18accc3971a43",

}



```
run

     ./pancakeswap-console

Update Price Every Second
![image](https://raw.githubusercontent.com/walletConsole/pancakeswap-console/master/image/6.jpg)



### Note
* The price may be a little different from that of the official website.Because using route path  <token> -> WBNB -> USDT
* Don't set up too many trading pairs, there will be problems with rpc node.


### Donations

BNB/CAKE/USDT: 0xec5fa25e37dfa8fa42210a94cbc8a61c7fd3751c

email: irmakgu40@gmail.com