package main

import (
	"github.com/nomadcoin/blockchain"
	"github.com/nomadcoin/cli"
)

func main() {
	blockchain.Blockchain()
	cli.Start()
}
