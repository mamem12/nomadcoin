package main

import (
	"github.com/nomadcoin/blockchain"
	"github.com/nomadcoin/cli"
	"github.com/nomadcoin/db"
)

func main() {
	defer db.Close()
	blockchain.Blockchain()
	cli.Start()
}
