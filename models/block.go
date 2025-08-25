package models

import "gorm.io/gorm"

type Block struct {
	gorm.Model
	Index        int
	Transactions []Transaction `gorm:"foreignKey:BlockID"`
	PrevHash     string
	Hash         string
	Timestamp    int64
	MerkleRoot   string
}
