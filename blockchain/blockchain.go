package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

type blockChain struct {
	blocks []*Block
}

func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{
		Data:     data,
		Hash:     "",
		PrevHash: getLastHash(),
		Height:   len(GetBlockchain().blocks) + 1,
	}

	newBlock.calculateHash()

	return &newBlock
}

var b *blockChain
var once sync.Once

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func (b *blockChain) Addblock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockchain() *blockChain {
	if b == nil {
		once.Do(func() {
			b = &blockChain{}
			b.Addblock("Genesis Block")
		})
	}
	return b
}

func (b *blockChain) AllBlocks() []*Block {
	return b.blocks
}

func (b *blockChain) GetBlock(height int) *Block {
	return b.blocks[height-1]
}
