package schedule

import (
	"blockchain/controllers"
	"blockchain/database"
	"blockchain/helpers"
	"blockchain/models"
	"fmt"
	"time"
)

var bc controllers.Blockchain

func AutoMineBlock() {
	var transactions []models.Transaction
	if result := database.DB.Find(&transactions); result.Error != nil {
		helpers.ErrorResponse(nil, result.Error.Error())
		return
	}

	if len(transactions) == 0 {
		fmt.Println("No transactions found")
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
}
