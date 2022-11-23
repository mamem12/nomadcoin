package main

import (
	"github.com/nomadcoin/cli"
	"github.com/nomadcoin/db"
)

func main() {
	defer db.Close()
	cli.Start()

}
