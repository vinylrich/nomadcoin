package main

import (
	"nomadcoin/explorer"
	"nomadcoin/rest"
)

func main() {
	go explorer.Start(3000)
	rest.Start(4000)
}
