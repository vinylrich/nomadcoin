package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const port string = ":3000"

type URLDes struct {
	URL     string `json:"url"`
	Method  string `json:"method"`
	Desc    string `json:"description"`
	Payload string `json:"payload,omitemipty"`
}

//uppercase로 하되 json으로 보내지는 이름을
//바꾸고 싶으면 json 태크를 사용해라
func documentation(w http.ResponseWriter, r *http.Request) {
	data := []URLDes{
		{
			URL:    "/",
			Method: "GET",
			Desc:   "See Documentation",
		},

		{
			URL:     "/block",
			Method:  "POST",
			Desc:    "Create Block",
			Payload: "data:string",
		},
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/", documentation)
	log.Printf("ListenAndServe http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
