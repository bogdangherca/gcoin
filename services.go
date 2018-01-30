package main

import (
	"net/http"
	"encoding/json"
	"log"
)

type NodeAddressList struct {
	addressList []string `json:"address_list"`
}

func newTransaction(w http.ResponseWriter, request *http.Request) {
	var transaction Transaction

	err := json.NewDecoder(request.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// TODO: validate non-empty fields
	blockchain.newTransaction(transaction.Sender, transaction.Recipient, transaction.Amount)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func mine(w http.ResponseWriter, request *http.Request) {
	reward := 1
	lastProof := blockchain.lastBlock().Proof

	// start mining
	log.Printf("node %v started mining block with last_proof = %v", nodeId, lastProof)
	proof := proofOfWork(lastProof)
	log.Printf("node %v finished mining block with proof = %v", nodeId, proof)

	// add reward transaction
	blockchain.newTransaction("0", nodeId, reward)

	// add new block with obtained proof
	minedBlock := blockchain.newBlock(proof, blockchain.lastBlock().hash())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(minedBlock)
}

func registerNode(w http.ResponseWriter, request *http.Request) {
	newNodes := NodeAddressList{}
	json.NewDecoder(request.Body).Decode(newNodes)

	for _, newNode := range newNodes.addressList {
		log.Printf("adding neighbour node %v to node id %v", newNode, nodeId)
		blockchain.Nodes = append(blockchain.Nodes, newNode)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}