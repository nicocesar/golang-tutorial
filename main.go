package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nicocesar/golang-tutorial/lib/contracts/erc20"
)

func main() {
	// Go get your key:  https://infura.io/register
	key := os.Getenv("INFURA_API_KEY")
	if key == "" {
		fmt.Println("INFURA_API_KEY is not set")
		os.Exit(1)
	}

	// Connect to Ethereum node
	client, err := ethclient.Dial(fmt.Sprintf("wss://mainnet.infura.io/ws/v3/%s", key))
	if err != nil {
		panic(err)
	}

	token0Address := common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F") // DAI mainnet
	token0contract, err := erc20.NewERC20(token0Address, client)
	if err != nil {
		panic(err)
	}

	for signature, methods := range erc20.ERC20MetaData.Sigs {
		fmt.Println(signature, methods)
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

	fmt.Printf("Contract %s, Symbol:%s, Name: %s, Decimals %d \n", token0Address, symbol, name, decimals)
}
