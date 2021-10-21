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
	port        string = ":3000"
	templateDir string = "explorer/templates/"
)

var templates *template.Template

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", blockchain.GetBlockchain().ListOfBlocks()}
	templates.ExecuteTemplate(w, "home", data)
}
func addHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(w, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("data")
		blockchain.GetBlockchain().AddBlock(data)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}

}

func Start() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
