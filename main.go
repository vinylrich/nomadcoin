package main

import (
	"fmt"
	"log"
	"net/http"
	"nomadcoin/blockchain"
	"text/template"
)

type homeData struct {
	PageTitle string //has to be uppercase
	Blocks    []*blockchain.Block
}

const port string = ":3000"

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	data := homeData{"Home", blockchain.GenerateBlockchain().ListOfBlocks()}
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", indexHandler)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
