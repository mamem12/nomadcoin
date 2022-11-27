package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nomadcoin/blockchain"
	"github.com/nomadcoin/utils"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type addTxPayload struct {
	To     string
	Amount int
}

func (u urlDescription) String() string {
	return "Hello I'm URL description"
}

func documentation(w http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the status of the Blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an address",
			Payload:     "data:string",
		},
	}
	json.NewEncoder(w).Encode((data))
}

func blocks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// return
		err := json.NewEncoder(w).Encode(blockchain.Blocks(blockchain.Blockchain()))
		utils.HandleErr(err)
	case "POST":
		blockchain.Blockchain().Addblock()
		w.WriteHeader(http.StatusCreated)
	}
}

func block(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	block, err := blockchain.FindBlock(hash)

	encoder := json.NewEncoder(w)

	if err == blockchain.ErrNotFound {
		utils.HandleErr(err)
	} else {
		utils.HandleErr(encoder.Encode(block))
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func status(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(blockchain.Blockchain())
}

func balance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.Blockchain())
		json.NewEncoder(w).Encode(balanceResponse{
			Address: address,
			Balance: amount,
		})
	default:
		err := json.NewEncoder(w).Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain()))
		utils.HandleErr(err)

	}
}

func mempool(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(blockchain.Mempool.Txs)
	utils.HandleErr(err)
}

func transactions(w http.ResponseWriter, r *http.Request) {
	payload := &addTxPayload{}
	err := json.NewDecoder(r.Body).Decode(payload)
	utils.HandleErr(err)
	err = blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		json.NewEncoder(w).Encode(errorResponse{"not enough funds"})
	}
	w.WriteHeader(http.StatusCreated)
}

func Start(aPort int) {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	fmt.Printf("listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
