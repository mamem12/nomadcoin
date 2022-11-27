package blockchain

import (
	"sync"

	"github.com/nomadcoin/db"
	"github.com/nomadcoin/utils"
)

type blockChain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

var b *blockChain
var once sync.Once

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval          = 2
	allowedRange           = 2
)

func (b *blockChain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockChain) Addblock() {
	block := CreateBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
}

func Blocks(b *blockChain) []*Block {
	hashCursor := b.NewestHash
	var blocks []*Block
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}

	return blocks
}

func persistBlockchain(b *blockChain) {
	db.SaveCheckpoint(utils.ToBytes(b))
}

func recalculateDifficulty(b *blockChain) int {
	allBlocks := Blocks(b)
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	// 분단위 계산
	expectedTime := difficultyInterval * blockInterval
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func getDifficulty(b *blockChain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

func UTxOutsByAddress(address string, b *blockChain) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					creatorTxs[input.TxID] = true
				}
			}

			for index, output := range tx.TxOuts {
				if output.Owner == address {
					if _, ok := creatorTxs[tx.ID]; !ok {
						uTxOut := &UTxOut{tx.ID, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}

				}

			}
		}
	}
	return uTxOuts
}

func BalanceByAddress(address string, b *blockChain) int {
	txOuts := UTxOutsByAddress(address, b)
	var amount int

	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

func Blockchain() *blockChain {
	once.Do(func() {
		b = &blockChain{
			NewestHash: "",
			Height:     0,
		}
		checkPoint := db.CheckPoint()
		if checkPoint == nil {
			b.Addblock()
		} else {
			b.restore(checkPoint)
		}
	})
	return b
}
