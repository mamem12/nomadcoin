package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct {
	Data     string
	Hash     string
	PrevHash string
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
		PrevHash: getLastHash(),
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
