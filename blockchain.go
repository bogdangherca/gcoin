package main

import "time"
import "encoding/json"
import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"log"
	"net/http"
	"fmt"
)

type Blockchain struct {
	Chain               []Block `json:"chain"`
	CurrentTransactions []Transaction `json:"current_transactions"`
	Nodes []string `json:"nodes"`
}

type Transaction struct {
	// complete here
	Sender string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount int `json:"amount"`
}

type Block struct {
	Index            int `json:"index"`
	Timestamp        time.Time `json:"timestamp"`
	Transaction_list []Transaction `json:"transaction_list"`
	Proof            int `json:"proof"`
	PrevHash         string `json:"prev_hash"`
}

// creates new block in the blockchain
func (b *Blockchain) newBlock(proof int, prev_hash string) Block {
	// prev_hash parameter is optional
	if prev_hash == "" {
		prev_hash = b.lastBlock().hash()
	}

	block := Block{
		len(b.Chain) + 1,
		time.Now(),
		b.CurrentTransactions,
		proof,
		prev_hash,
	}

	// reset the current list of transactions
	b.CurrentTransactions = []Transaction{}

	b.Chain = append(b.Chain, block)

	return block
}

// creates a new transaction to go into the next mined Block
func (b *Blockchain) newTransaction(sender string, recipient string, amount int) int {
	transaction := Transaction{
		sender,
		recipient,
		amount,
	}

	b.CurrentTransactions = append(b.CurrentTransactions, transaction)

	return b.lastBlock().Index + 1
}

// fetch the last block
func (b *Blockchain) lastBlock() Block {
	return b.Chain[len(b.Chain) - 1]
}

// creates a SHA-256 hash of a Block
func (block Block) hash() string {
	byteBlock, _ := json.Marshal(block)
	shaBlock := sha256.Sum256(byteBlock)
	return hex.EncodeToString(shaBlock[:])
}

/*
Simple Proof of Work Algorithm:
   - Find a number p' such that hash(p*p') contains leading 4 zeroes, where p is the previous p'
   - p is the previous proof, and p' is the new proof
*/
func proofOfWork(state int) int {

	proof := 0

	for validateProofOfWork(state, proof) == false {
		proof++
	}

	return proof
}

func validateProofOfWork(lastProof int, proof int) bool {
	// modify this to increase the difficulty
	difficulty := 5
	leadingZero := strings.Repeat("0", difficulty)

	byteProof, _ := json.Marshal(fmt.Sprintf("%v%v", lastProof, proof))
	shaProof := sha256.Sum256(byteProof)
	strProof := hex.EncodeToString(shaProof[:])

	if strings.HasPrefix(strProof, leadingZero) == true {
		log.Printf("proof checking passed for proof %v with hash %v", proof, strProof)
	}

	return strings.HasPrefix(strProof, leadingZero)
}

func (b *Blockchain) validateChain() bool {
	for i, block := range b.Chain[1:] {
		prevBlock := b.Chain[i-1]
		if block.PrevHash != prevBlock.hash() {
			return false
		}

		if validateProofOfWork(prevBlock.Proof, block.Proof) == false {
			return false
		}
	}

	return true
}

func (b *Blockchain) resolveConflicts() bool {
	for _, node := range b.Nodes {
		var otherBlockchain Blockchain
		resp, _ := http.Get(node + "/chain")
		defer resp.Body.Close()

		log.Printf("resp.StatusCode=%v", resp.StatusCode)
		if resp.StatusCode == 200 {
			json.NewDecoder(resp.Body).Decode(otherBlockchain)

			// replace current chain only if we found a longer valid one
			log.Printf("len(otherBlockchain.Chain) > len(b.Chain) = %v", len(otherBlockchain.Chain) > len(b.Chain))
			log.Printf("otherBlockchain.validateChain() = %v", otherBlockchain.validateChain())
			if len(otherBlockchain.Chain) > len(b.Chain) && otherBlockchain.validateChain() == true {
				b.Chain = otherBlockchain.Chain
				return true
			}
		}
	}

	return false
}

// register a new node by address eg. 'http://192.168.0.5:5000'
func (b *Blockchain) registerNode(address string) {
	b.Nodes = append(b.Nodes, address)
}
