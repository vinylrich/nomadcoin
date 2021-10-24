package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const port string = ":3000"

type URL string

//go에서는 상속,implement가 없기 때문에
//reseiver를 사용해서 명시 없이 implement한다
//아래와 같이 구현하면 url이 marshaltext,
//urldec가 string을 implement한 것이다.
func (u URL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type URLDes struct {
	URL     URL    `json:"url"`
	Method  string `json:"method"`
	Desc    string `json:"description"`
	Payload string `json:"payload,omitempty"`
}

func (u URLDes) String() string {
	return "Hello Im the URLDEC"
}

//uppercase로 하되 json으로 보내지는 이름을
//바꾸고 싶으면 json 태그를 사용해라
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
