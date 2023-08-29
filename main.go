package main

import (
	"fmt"

	blockchain "github.com/tensor-programming/golang-blockchain/Blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()
	chain.AddBlock("Primeiro bloco depois do Gênesis!")
	chain.AddBlock("Segundo bloco depois do Gênesis!")
	chain.AddBlock("Terceiro bloco depois do Gênesis!")

	for _, block := range chain.Blocks {
		fmt.Printf("Hash anterior: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}
}
