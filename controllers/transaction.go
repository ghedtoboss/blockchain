package controllers

import (
	"blockchain/database"
	"blockchain/helpers"
	"blockchain/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func TransactionToString(tx models.Transaction) string {
	return fmt.Sprintf("%s->%s:%s:%.2f:%.2f:%s",
		tx.From, tx.To, tx.Currency, tx.Amount, tx.Fee, tx.Signature)
}

func TransactionsToString(transactions []models.Transaction) []string {
	var strings []string
	for _, value := range transactions {
		strings = append(strings, TransactionToString(value))
	}
	return strings
}

func CalculateSHA256FromTransactions(transactions []models.Transaction) string {
	strings := TransactionsToString(transactions)
	return CalculateSHA256(strings)
}

func CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		helpers.ErrorResponse(c, err.Error())
		return
	}

	if result := database.DB.Create(&transaction); result.Error != nil {
		helpers.ErrorResponse(c, result.Error.Error())
		return
	}

	helpers.SuccessResponse(c, transaction)
}
