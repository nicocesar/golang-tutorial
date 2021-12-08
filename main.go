package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nicocesar/golang-tutorial/lib/contracts/erc20"
	"github.com/nicocesar/golang-tutorial/lib/contracts/uniswapv3"
)

func describeERC20(tokenAddress common.Address, client *ethclient.Client) (string, error) {
	tokenContract, err := erc20.NewERC20(tokenAddress, client)
	if err != nil {
		panic(err)
	}

	name, err := tokenContract.Name(nil)
	if err != nil {
		return "", err
	}

	symbol, err := tokenContract.Symbol(nil)
	if err != nil {
		return "", err
	}

	decimals, err := tokenContract.Decimals(nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Token: %s, Symbol: %s, Name: %s, Decimals: %d\n", tokenAddress, symbol, name, decimals), nil
}

func main() {
	// Go get your key:  https://infura.io/register
	key := os.Getenv("INFURA_API_KEY")
	if key == "" {
		fmt.Println("INFURA_API_KEY is not set. Get an account at https://infura.io/register")
		os.Exit(1)
	}

	// Connect to Ethereum node
	client, err := ethclient.Dial(fmt.Sprintf("wss://mainnet.infura.io/ws/v3/%s", key))
	if err != nil {
		panic(err)
	}

	poolAddress := common.HexToAddress("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640") // ETH-USDC 0.05% UniV3 Pool
	boundcontract, err := uniswapv3.NewUniswapV3Pool(poolAddress, client)
	if err != nil {
		panic(err)
	}
	c := boundcontract.UniswapV3PoolCaller

	for signature, methods := range uniswapv3.IUniswapV3PoolMetaData.Sigs {
		fmt.Println(signature, methods)
	}

	token0Address, err := c.Token0(nil)
	if err != nil {
		panic(err)
	}
	tokenStr, err := describeERC20(token0Address, client)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Token 0: %s", tokenStr)

	token1Address, err := c.Token1(nil)
	if err != nil {
		panic(err)
	}
	tokenStr, err = describeERC20(token1Address, client)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Token 1: %s", tokenStr)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{poolAddress},
	}
	u := boundcontract.UniswapV3PoolFilterer

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:

			fmt.Printf("https://etherscan.io/tx/%s\n", vLog.TxHash.Hex())

			switch vLog.Topics[0] {
			// Swap
			case crypto.Keccak256Hash([]byte("Swap(address,address,int256,int256,uint160,uint128,int24)")):
				y, err := u.ParseSwap(vLog)
				if err != nil {
					log.Printf("error:%s\n%#v", err.Error(), vLog)
					break
				}

				fmt.Printf("Swap: %d,%d,%d,%d,%d\n", y.Amount0, y.Amount1, y.Liquidity, y.Tick, y.SqrtPriceX96)
			default:
				fmt.Printf("Unknown log: %#v\n", vLog)
			}
		}
	}
}
