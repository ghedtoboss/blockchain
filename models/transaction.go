package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	BlockID   uint
	From      string
	To        string
	Currency  string
	Amount    float64
	Fee       float64
	Signature string
}
