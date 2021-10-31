package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nomadcoin/blockchain"
	"nomadcoin/utils"

	"github.com/gorilla/mux"
)

var port string

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}
type url string

type blockBody struct {
	Message string `json:"data"`
}

//go에서는 상속,implement가 없기 때문에
//reseiver를 사용해서 명시 없이 implement한다
//아래와 같이 구현하면 url이 marshaltext,
//urldec가 string을 implement한 것이다.
func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type uRLDes struct {
	URL     url    `json:"url"`
	Method  string `json:"method"`
	Desc    string `json:"description"`
	Payload string `json:"payload,omitempty"`
}

//uppercase로 하되 json으로 보내지는 이름을
//바꾸고 싶으면 json 태그를 사용해라
func documentation(w http.ResponseWriter, r *http.Request) {
	data := []uRLDes{
		{
			URL:    "/",
			Method: "GET",
			Desc:   "See Documentation",
		},

		{
			URL:     "/blocks",
			Method:  "POST",
			Desc:    "Create Block",
			Payload: "data:string",
		},
		{
			URL:    "/blocks",
			Method: "GET",
			Desc:   "GET ALL Block",
		},
		{
			URL:    "/blocks/{height}",
			Method: "GET",
			Desc:   "GET Specific Block",
		},
	}
	json.NewEncoder(w).Encode(data)
}

func getAllBlocks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(blockchain.GetBlockchain().ListOfBlocks())
}

func createBlock(w http.ResponseWriter, r *http.Request) {
	var blockBody blockBody
	utils.HandleError(json.NewDecoder(r.Body).Decode(&blockBody))

	blockchain.GetBlockchain().AddBlock(blockBody.Message)
	w.WriteHeader(http.StatusCreated)
}
func getBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(w)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	}
	encoder.Encode(block)
}
func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", getAllBlocks).Methods("GET")
	router.HandleFunc("/blocks", createBlock).Methods("POST")
	router.HandleFunc("/blocks/{hash:[a-f]+}", getBlock).Methods("GET")
	log.Printf("ListenAndServe http://localhost%s to rest api\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
