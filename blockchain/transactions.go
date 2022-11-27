package blockchain

import (
	"errors"
	"fmt"
	"time"

	"github.com/nomadcoin/utils"
)

const (
	minerReward int = 50
)

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type TxIn struct {
	TxID  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{
			TxID:  "",
			Index: -1,
			Owner: "COINBASE",
		},
	}
	txOuts := []*TxOut{
		{
			Owner:  address,
			Amount: minerReward,
		},
	}
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx
}

func isOnMempool(uTxOut *UTxOut) bool {

	exists := false
Outer:
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			exists = input.TxID == uTxOut.TxID && input.Index == uTxOut.Index
			break Outer
		}
	}

	return exists
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, errors.New("not enoguh 돈")
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutsByAddress(from, Blockchain())
	for _, uTxOut := range uTxOuts {
		if total > amount {
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	if change := total - amount; change != 0 {
		fmt.Printf("Change :: %d\n", change)
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {

	tx, err := makeTx("MAMEM", to, amount)
	if err != nil {
		return err
	}

	m.Txs = append(m.Txs, tx)

	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("MAMEM")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
