package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nomadcoin/utils"
)

const port string = ":3000"

type URLDes struct {
	URL    string
	Method string
	Desc   string
}

func documentation(w http.ResponseWriter, r *http.Request) {
	data := []URLDes{
		{
			URL:    "/",
			Method: "GET",
			Desc:   "See Documentation",
		},

		{
			URL:    "/block",
			Method: "POST",
			Desc:   "Create Block",
		},
	}
	pbyte, err := json.Marshal(data)
	utils.HandleError(err)
	fmt.Fprint(w, string(pbyte))
}

func main() {
	http.HandleFunc("/", documentation)
	log.Printf("ListenAndServe http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
