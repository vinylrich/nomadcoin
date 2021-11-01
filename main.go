package main

import (
	"nomadcoin/blockchain"
	"nomadcoin/cli"
	"nomadcoin/db"
)

func main() {
	defer db.Close()
	blockchain.Blockchain()
	cli.Start()
}
