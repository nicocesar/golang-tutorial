package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/nicocesar/golang-tutorial/lib/contracts/erc20"
	"github.com/nicocesar/golang-tutorial/lib/contracts/uniswapv3"
	"github.com/shopspring/decimal"
)

func describeERC20(tokenAddress common.Address, client *ethclient.Client) (string, string, uint8, error) {
	tokenContract, err := erc20.NewERC20(tokenAddress, client)
	if err != nil {
		return "", "", 0, err
	}

	name, err := tokenContract.Name(nil)
	if err != nil {
		return "", "", 0, err
	}

	symbol, err := tokenContract.Symbol(nil)
	if err != nil {
		return "", "", 0, err
	}

	decimals, err := tokenContract.Decimals(nil)
	if err != nil {
		return "", "", 0, err
	}

	return symbol, name, decimals, nil
}

func bigIntToDecimalString(amount *big.Int, decimals uint8) (string, error) {
	if decimals == 0 {
		return amount.String(), nil
	} else {
		return decimal.NewFromBigInt(amount, -int32(decimals)).String(), nil
	}
}

func prettySwap(Amount0 *big.Int, Amount1 *big.Int, t0Symbol string, t1Symbol string, t0Decimals uint8, t1Decimals uint8) (string, error) {
	var err error
	var t0Amount, t1Amount string
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintfFunc()

	sellingT0 := Amount0.Cmp(big.NewInt(0)) == -1 //  Amount0 is lower than 0 means that we are selling token0

	if t0Decimals > 0 {
		t0Amount, err = bigIntToDecimalString(Amount0, t0Decimals)
		if err != nil {
			return "", err
		}
	} else {
		t0Amount = fmt.Sprintf("%d", Amount0)
	}

	if t1Decimals > 0 {
		t1Amount, err = bigIntToDecimalString(Amount1, t1Decimals)
		if err != nil {
			return "", err
		}
	} else {
		t1Amount = fmt.Sprintf("%d", Amount1)
	}

	if sellingT0 {
		return fmt.Sprintf("%s %s -> %s %s", red(t0Amount), red(t0Symbol), green(t1Amount), green(t1Symbol)), nil
	} else {
		return fmt.Sprintf("%s %s -> %s %s", green(t0Amount), green(t0Symbol), red(t1Amount), red(t1Symbol)), nil
	}

}

func showSwaps(poolAddress common.Address, client *ethclient.Client) {
	fmt.Println("Watching swaps for pool: ", poolAddress.String())
	boundcontract, err := uniswapv3.NewUniswapV3Pool(poolAddress, client)
	if err != nil {
		panic(err)
	}
	c := boundcontract.UniswapV3PoolCaller

	token0Address, err := c.Token0(nil)
	if err != nil {
		panic(err)
	}
	t0Symbol, t0Name, t0Decimals, err := describeERC20(token0Address, client)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Token 0: %s (%s), Decimals: %d\n", t0Symbol, t0Name, t0Decimals)

	token1Address, err := c.Token1(nil)
	if err != nil {
		panic(err)
	}
	t1Symbol, t1Name, t1Decimals, err := describeERC20(token1Address, client)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Token 1: %s (%s), Decimals: %d\n", t1Symbol, t1Name, t1Decimals)

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
				x, err := prettySwap(y.Amount0, y.Amount1, t0Symbol, t1Symbol, t0Decimals, t1Decimals)
				if err != nil {
					log.Printf("error: %s\n%#v", err.Error(), vLog)
					break
				}
				fmt.Println(x)
			default:
				fmt.Printf("Unknown log: %#v\n", vLog)
			}
		}
	}

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

	pools, err := getUniswapV3Pools()
	if err != nil {
		panic(err)
	}
	poolsToWatch := make(chan common.Address, 10)

	go func() {
		// this will push the pool addresses to the channel to watch
		// this is a channel so in the future we can add more pools
		// dynamically from other go routines
		for _, pool := range pools.Pools {
			poolsToWatch <- common.HexToAddress(pool)
		}
	}()

	for poolAddress := range poolsToWatch {
		go showSwaps(poolAddress, client)
	}
}
