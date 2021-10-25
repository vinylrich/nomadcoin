package cli

import (
	"flag"
	"fmt"
	"nomadcoin/explorer"
	"nomadcoin/rest"
	"os"
)

func intro() {
	fmt.Printf("Welcome to NOMAD COIN\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port=4000: 	Set port of the server \n")
	fmt.Printf("-mode=rest: 	Choose between 'html and 'rest'\n")
	os.Exit(0)
}

func Start() {
	if len(os.Args) == 1 {
		intro()
	}
	port := flag.Int("port", 4000, "Sets the port of the server")

	mode := flag.String("mode", "rest", "Choose between 'html and 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		intro()
	}
	fmt.Println(*port, *mode)
}

// command는 go run main.go explorer
// go run main.go rest
//flag는 -port 4000
