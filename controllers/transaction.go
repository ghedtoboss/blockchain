package controllers

import (
	"blockchain/models"
	"fmt"
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
