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

const (
	port        string = ":3000"
	templateDir string = "templates/"
)

var templates *template.Template

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", blockchain.GetBlockchain().ListOfBlocks()}
	templates.ExecuteTemplate(w, "home", data)
}

func main() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	http.HandleFunc("/", indexHandler)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
