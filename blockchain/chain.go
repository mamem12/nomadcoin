package blockchain

import (
	"fmt"
	"sync"

	"github.com/nomadcoin/db"
	"github.com/nomadcoin/utils"
)

type blockChain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockChain
var once sync.Once

func (b *blockChain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockChain) persist() {
	db.SaveBlockChain(utils.ToBytes(b))
}

func (b *blockChain) Addblock(data string) {
	block := CreateBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func Blockchain() *blockChain {
	if b == nil {
		once.Do(func() {
			b = &blockChain{
				NewestHash: "",
				Height:     0,
			}
			checkPoint := db.CheckPoint()
			if checkPoint == nil {
				b.Addblock("Genesis Block")
			} else {
				b.restore(checkPoint)
			}
		})
	}
	fmt.Printf("Newesthash : %s\n", b.NewestHash)
	return b
}
