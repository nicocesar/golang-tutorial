package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nicocesar/golang-tutorial/lib/contracts/erc20"
	"github.com/nicocesar/golang-tutorial/lib/contracts/uniswapv3"
)

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
	token0contract, err := erc20.NewERC20(token0Address, client)
	if err != nil {
		panic(err)
	}

	name, err := token0contract.ERC20Caller.Name(nil)
	if err != nil {
		panic(err)
	}

	symbol, err := token0contract.Symbol(nil)
	if err != nil {
		panic(err)
	}

	decimals, err := token0contract.Decimals(nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Token0: %s, Symbol: %s, Name: %s, Decimals: %d\n", token0Address, symbol, name, decimals)
}
