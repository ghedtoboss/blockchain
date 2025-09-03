package controllers

import (
	"blockchain/database"
	"blockchain/helpers"
	"blockchain/models"
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Blockchain struct {
	Blocks              []models.Block
	PendingTransactions []models.Transaction
}

var bc Blockchain

func (bc *Blockchain) SaveBlockDB(block models.Block) {
	if result := database.DB.Create(&block); result.Error != nil {
		fmt.Println("Error saving block to database: ", result.Error)
	}
}

func (bc *Blockchain) CreateBlock(prevBlock models.Block, transactions []models.Transaction, timestamp int64, difficulty int) models.Block {
	var block models.Block
	block.Index = prevBlock.Index + 1
	block.PrevHash = prevBlock.Hash
	block.Transactions = transactions
	block.Timestamp = timestamp
	block.MerkleRoot = CalculateSHA256FromTransactions(transactions)

	bc.MineBlock(&block, difficulty)
	bc.Blocks = append(bc.Blocks, block)
	bc.SaveBlockDB(block)

	return block
}

func CalculateBlockHash(block models.Block) string {
	return CalculateSHA256([]string{
		block.PrevHash,
		block.MerkleRoot,
		fmt.Sprintf("%d", block.Timestamp),
		fmt.Sprintf("%d", block.Nonce)})
}

func (bc *Blockchain) CreateGenesisBlock(timestamp int64) models.Block {
	var block models.Block
	block.Index = 0
	block.PrevHash = "0000"
	block.Transactions = []models.Transaction{
		{
			From:      "0",
			To:        "0",
			Currency:  "0",
			Amount:    0,
			Fee:       0,
			Signature: "0",
		},
	}
	block.Timestamp = timestamp
	block.MerkleRoot = CalculateSHA256FromTransactions(block.Transactions)

	bc.MineBlock(&block, 3)
	bc.Blocks = append(bc.Blocks, block)
	bc.SaveBlockDB(block)
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

func (bc *Blockchain) ValidateBlock(prevBlock, currentBlock models.Block) bool {
	// index kontrolü
	if currentBlock.Index != prevBlock.Index+1 {
		return false
	}

	// prev hash kontrol
	if currentBlock.PrevHash != prevBlock.Hash {
		return false
	}

	// hash kontrol
	if currentBlock.Hash != CalculateBlockHash(currentBlock) {
		return false
	}

	// merkle root kontrol
	if currentBlock.MerkleRoot != CalculateSHA256FromTransactions(currentBlock.Transactions) {
		return false
	}

	return true
}

func (bc *Blockchain) ValidateChain() bool {
	// Genesis block kontrol
	if bc.Blocks[0].Index != 0 {
		return false
	}

	//  Her komşu block çifti için validateblock
	for i := 1; i < len(bc.Blocks); i++ {
		if !bc.ValidateBlock(bc.Blocks[i-1], bc.Blocks[i]) {
			return false
		}
	}
	return true
}

func (bc *Blockchain) MineBlock(block *models.Block, difficulty int) {
	// Target: "000.." (difficulty kadar sıfır)
	target := strings.Repeat("0", difficulty)

	// uygun hash bulana kadar dene
	for {
		hash := CalculateBlockHash(*block)

		// Hash target ile uyuşuyor mu
		if strings.HasPrefix(hash, target) {
			block.Hash = hash
			break // bulundu
		}

		// uymuyor, nonce'u arttır ve tekrar dene
		block.Nonce++
	}
}

// tüm blockchaini JSON olarak döndür
func GetChain(c *gin.Context) {

	var blocks []models.Block
	database.DB.Find(&blocks)

	c.JSON(http.StatusOK, blocks)
}

func ManuelMineBlock(c *gin.Context) {
	var transactions []models.Transaction
	if result := database.DB.Find(&transactions); result.Error != nil {
		helpers.ErrorResponse(c, result.Error.Error())
		return
	}

	if len(transactions) == 0 {
		helpers.ErrorResponse(c, "No transactions found")
		return
	}

	var prevBlock models.Block
	if result := database.DB.Last(&prevBlock); result.Error != nil {
		var blockCount int64
		database.DB.Model(&models.Block{}).Count(&blockCount)
		if blockCount == 0 {
			bc.CreateGenesisBlock(time.Now().Unix())
		} else {
			fmt.Println("Previous block found")
			return
		}
	}

	fmt.Printf("Mining block with %d transactions...\n", len(transactions))
	bc.CreateBlock(prevBlock, transactions, time.Now().Unix(), 3)

	database.DB.Delete(&transactions)
	fmt.Println("Block mined")

	helpers.SuccessResponse(c, "Block mined")
}
