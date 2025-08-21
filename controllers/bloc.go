package controllers

import (
	"blockchain/models"
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

func CreateBlock(prevHash string, index int,transactions []string) models.Block {
	var block models.Block
	block.Index = index
	block.PrevHash = prevHash // Önceki bloğun hash'i
	block.Transactions = transactions
	block.Timestamp = time.Now().Unix()
	block.MerkleRoot = CalculateSHA256(transactions)
	blockHash := CalculateBlockHash(block) // Yeni hash hesaplanıyor
	block.Hash = blockHash
	return block
}

func CalculateBlockHash(block models.Block) string {
	return CalculateSHA256([]string{block.PrevHash, block.MerkleRoot, fmt.Sprintf("%d", block.Timestamp)})
}

func CreateGenesisBlock() models.Block {
	var block models.Block
	block.Index = 0
	block.PrevHash = "0000" // İlk bloğun hash'i
	block.Transactions = []string{"Ali -> Veli 10 BTC"}
	block.Timestamp = time.Now().Unix()
	block.MerkleRoot = CalculateSHA256([]string{"Ali -> Veli 10 BTC"})
	blockHash := CalculateBlockHash(block) // Yeni hash hesaplanıyor
	block.Hash = blockHash
	return block
}

func CalculateSHA256(data []string) string {
	hashes := []string{}

	for _, value := range data {
		dataBytes := []byte(value)
		hash := sha256.Sum256(dataBytes)
		hashes = append(hashes, fmt.Sprintf("%x", hash))
	}

	for len(hashes) > 1 {
		newHashes := []string{} // yeni seviye için boş slice

		for i := 0; i < len(hashes); i += 2 {
			left := hashes[i]
			if i+1 < len(hashes) {
				right := hashes[i+1]
				newHashes = append(newHashes, CombineHashes(left, right))
			} else {
				newHashes = append(newHashes, left)
			}

		}
		hashes = newHashes
	}
	return hashes[0]
}

func CombineHashes(left, right string) string {
	var buffer bytes.Buffer
	buffer.WriteString(left)
	buffer.WriteString(right)
	return fmt.Sprintf("%x", sha256.Sum256(buffer.Bytes()))
}
