package explorer

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
	templateDir string = "explorer/templates/"
)

var templates *template.Template

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", blockchain.Blockchain().AllBlocks()}
	templates.ExecuteTemplate(w, "home", data)
}
func addHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(w, "add", nil)
	case "POST":
		r.ParseForm()
		blockchain.Blockchain().AddBlock()
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

}

func Start(port int) {
	handler := http.NewServeMux()
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	handler.HandleFunc("/", indexHandler)
	handler.HandleFunc("/add", addHandler)
	fmt.Printf("Listening on http://localhost:%d to html\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
