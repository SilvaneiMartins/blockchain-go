package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *Blockchain {
	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("Blockchain não encontrado. Criando um novo...")

			genesis := Genesis()
			println("Genesis criado")

			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)

			err = txn.Set([]byte("lh"), genesis.Hash)
			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)

			err1 := item.Value(func(val []byte) error {
				lastHash = append([]byte{}, val...)
				return nil
			})
			Handle(err1)

			return err
		}
	})
	Handle(err)
	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		err1 := item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})

		return err1
	})

	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)

		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash
		return err
	})

	Handle(err)
}

func (chain *Blockchain) Iterator() *BlockchainIterator {
	iter := &BlockchainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockchainIterator) Next() *Block {
	var block *Block
	var blockData []byte

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)

		err1 := item.Value(func(val []byte) error {
			blockData = append([]byte{}, val...)
			return nil
		})

		Handle(err1)

		block = Deserialize(blockData)
		return err
	})

	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}
