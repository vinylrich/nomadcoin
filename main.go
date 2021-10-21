package main

import (
	"fmt"
	"log"
	"net/http"
)

const port string = ":3000"

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello index")
}
func main() {
	http.HandleFunc("/", indexHandler)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
