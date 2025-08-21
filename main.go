package main

import (
	"blockchain/controllers"
	"fmt"
)

func main() {
	// 1. Genesis Block
	genesis := controllers.CreateGenesisBlock()
	fmt.Println("Genesis Hash: ", genesis.Hash)

	// 2. Block 1
	block1 := controllers.CreateBlock(genesis.Hash, genesis.Index+1, []string{"Veli -> Ayşe 5 BTC"})
	fmt.Println("Block 1 Hash: ", block1.Hash)

	// 3. Block 2
	block2 := controllers.CreateBlock(block1.Hash, block1.Index+1, []string{"Ayşe -> Ali 10 BTC"})
	fmt.Println("Block 2 Hash: ", block2.Hash)
}
