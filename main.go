package main

import (
	"fmt"

	"github.com/nicocesar/golang-tutorial/lib/contracts/erc20"
)

func main() {
	fmt.Println("Signatures of methods from an ERC20 Contract: ")
	for signature, methods := range erc20.ERC20MetaData.Sigs {
		fmt.Println(signature, methods)
	}
}
