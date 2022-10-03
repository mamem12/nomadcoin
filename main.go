package main

import (
	"fmt"

	"main.go/blockchain"
)

func main() {
	chain := blockchain.GetBlockchain()
	chain.Addblock("second block")
	chain.Addblock("third block")
	chain.Addblock("fourth block")

	for _, v := range chain.AllBlocks() {
		fmt.Printf("Data : %s\n", v.Data)
		fmt.Printf("Hash : %s\n", v.Hash)
		fmt.Printf("PrevHash : %s\n", v.PrevHash)
	}
}
