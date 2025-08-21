package models

type Block struct {
	Index        int
	Transactions []string
	PrevHash     string
	Hash         string
	Timestamp    int64
	MerkleRoot   string
}
