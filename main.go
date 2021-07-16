package main

import (
	"encoding/json"
	"fmt"
	"gochain/blockchain"
	"gochain/blockchain/wallet"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type block struct {
	Address string `json:"address"`
}

type chainhash struct {
	Prev []byte `json:"prev_hash"`
	Hash []byte `json:"hash"`
	Pow  string `json:"pow"`
}

type formatresponse struct {
	Result []string `json:"result"`
}

func createBlockChain(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	req := block{}

	_ = json.Unmarshal(body, &req)

	newChain := blockchain.InitBlockChain(req.Address)
	newChain.Database.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result": "create block success"})
}

func getWallet(w http.ResponseWriter, r *http.Request) {

	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"result": address})
	return
}

func printChain(w http.ResponseWriter, r *http.Request) {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iterator := chain.Iterator()

	for {
		block := iterator.Next()
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		// This works because the Genesis block has no PrevHash to point to.
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func listAdress_Wallet(w http.ResponseWriter, r *http.Request) {

	wallets, _ := wallet.CreateWallets()
	addressess := wallets.GetAllAddresses()
	var result []string
	for _, address := range addressess {
		result = append(result, address)
	}
	format := formatresponse{}
	format.Result = result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(format)
	return
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/create_blockchain", createBlockChain).Methods("POST")
	r.HandleFunc("/print_blockchain", printChain).Methods("GET")
	r.HandleFunc("/create_wallet", getWallet).Methods("GET")
	r.HandleFunc("/wallet_collec", listAdress_Wallet).Methods("GET")

	http.ListenAndServe(":8081", r)
}
