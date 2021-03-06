package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nomadcoin/blockchain"
	"nomadcoin/p2p"
	"nomadcoin/utils"
	"nomadcoin/wallet"

	"github.com/gorilla/mux"
)

var port string

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}
type url string

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type addTxPayload struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type myWalletResponse struct {
	Address string `json:"address"`
}

type addPeerPayload struct {
	Address string
	Port    string
}

//go에서는 상속,implement가 없기 때문에
//reseiver를 사용해서 명시 없이 implement한다
//아래와 같이 구현하면 url이 marshaltext,
//urldes가 string을 implement한 것이다.
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
			URL:    url("/"),
			Method: "GET",
			Desc:   "See Documentation",
		},
		{
			URL:    url("/status"),
			Method: "GET",
			Desc:   "See the Status of the Blockchain",
		},
		{
			URL:     url("/blocks"),
			Method:  "POST",
			Desc:    "Create Block",
			Payload: "data:string",
		},
		{
			URL:    url("/blocks"),
			Method: "GET",
			Desc:   "GET ALL Block",
		},
		{
			URL:    url("/blocks/{hash}"),
			Method: "GET",
			Desc:   "GET Specific Block",
		},
		{
			URL:    url("/balance/{address}"),
			Method: "GET",
			Desc:   "Get TxOuts for an Address",
		},
		{
			URL:    url("/ws"),
			Method: "GET",
			Desc:   "Upgrade to WebSockets",
		},
	}
	json.NewEncoder(w).Encode(data)
}

func getAllBlocks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(blockchain.AllBlocks(blockchain.Blockchain()))
}

func createBlock(w http.ResponseWriter, r *http.Request) {
	newBlock := blockchain.Blockchain().AddBlock()
	p2p.BroadcastNewBlock(newBlock)
	w.WriteHeader(http.StatusCreated)
}
func getBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	utils.HandleError(err)
	encoder := json.NewEncoder(w)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	}
	encoder.Encode(block)
}

//Get Blockchain data
func blockchainStatus(w http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.Blockchain(), w)
}

func balance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	if address == "" {
		log.Println("no address")
		return
	}
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.Blockchain())
		json.NewEncoder(w).Encode(balanceResponse{address, amount})
	default:
		txOut := blockchain.UTxOutsByAddress(address, blockchain.Blockchain())
		utils.HandleError(json.NewEncoder(w).Encode(txOut))
	}

}

func mempool(w http.ResponseWriter, r *http.Request) {
	utils.HandleError(json.NewEncoder(w).Encode(blockchain.Mempool().Txs))
}

func transaction(w http.ResponseWriter, r *http.Request) {
	payload := &addTxPayload{}
	utils.HandleError(json.NewDecoder(r.Body).Decode(&payload))
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{ErrorMessage: err.Error()})
		return
	}
	p2p.BroadcastNewTx(tx)
	w.WriteHeader(http.StatusCreated)
}

func myWallet(w http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(w).Encode(myWalletResponse{Address: address})
}

func peers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, port[1:], true)
		w.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(w).Encode(p2p.AllPeers(&p2p.Peers))
	}

}
func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleware, loggerMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blockchain", blockchainStatus).Methods("GET")
	router.HandleFunc("/blocks", getAllBlocks).Methods("GET")
	router.HandleFunc("/blocks", createBlock).Methods("POST")
	router.HandleFunc("/blocks/{hash:[a-f]+}", getBlock).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	router.HandleFunc("/transaction", transaction).Methods("POST")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peers", peers).Methods("GET", "POST")
	log.Printf("ListenAndServe http://localhost%s to rest api\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(w, r)
	})
}
