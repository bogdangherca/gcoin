package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"github.com/google/uuid"
)

// instantiate the blockchain
var blockchain = Blockchain{}

var nodeId = uuid.Must(uuid.NewRandom()).String()

func main() {

	// add the genesis block
	blockchain.newBlock(100, "genesis")

	router := mux.NewRouter()

	router.HandleFunc("/mine", mine).Methods("GET")

	router.HandleFunc("/transactions", newTransaction).Methods("PUT")

	router.HandleFunc("/chain", func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(blockchain)
	}).Methods("GET")

	router.HandleFunc("/nodes/register", registerNode).Methods("POST")

	router.HandleFunc("/nodes/refresh", func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		isAuthoritative := blockchain.resolveConflicts()
		log.Printf("node %v has an authoritative chain: %v", nodeId, isAuthoritative)
		json.NewEncoder(w).Encode(isAuthoritative)
	}).Methods("GET")

	http.Handle("/", router)

	log.Printf("Starting blockchain node with id %v on port 5000...", nodeId)
	log.Fatal(http.ListenAndServe(":5001", router))
}
