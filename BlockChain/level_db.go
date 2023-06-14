package BlockChain

import (
	"BlockWallet/Block"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"math/big"
)

func SaveBlockChainToLevelDB(bc *Blockchain) error {
	// 打开或创建 LevelDB 数据库
	db, err := leveldb.OpenFile("blockChain.db", nil)
	if err != nil {
		return err
	}
	defer db.Close()

	// 将区块数据转换为字节数组
	blockchainData, err := bc.MarshalJSON()
	if err != nil {
		return err
	}

	// 保存区块到 LevelDB

	err = db.Put([]byte("yzg blockchain"), blockchainData, nil)
	if err != nil {
		return err
	}

	return nil
}

func SaveBlockToLevelDB(b *Block.Block) error {
	// 打开或创建 LevelDB 数据库
	db, err := leveldb.OpenFile("blocks.db", nil)
	if err != nil {
		return err
	}
	defer db.Close()

	// 将区块数据转换为字节数组
	blockData, err := b.MarshalJSON()
	if err != nil {
		return err
	}

	// 保存区块到 LevelDB
	var tmpNum *big.Int
	tmpNum = b.Number()

	err = db.Put(tmpNum.Bytes(), blockData, nil)
	if err != nil {
		return err
	}

	return nil
}

func LoadBlockChainFromLevelDB() (*Blockchain, error) {
	db, err := leveldb.OpenFile("blockChain.db", nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	blockchainByte, err := db.Get([]byte("yzg blockchain"), nil)
	if err == leveldb.ErrNotFound {
		log.Println("没有能加载的区块链")
		return nil, err
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var tmpBlockchain Blockchain

	err = tmpBlockchain.UnmarshalJSON(blockchainByte)

	return &tmpBlockchain, err
}

func LoadBlockFromLevelDB(tmpBlockchain *Blockchain) (*Blockchain, error) {
	// 打开 LevelDB 数据库
	db, err := leveldb.OpenFile("blocks.db", nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 遍历 LevelDB 中的所有键值对
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// 获取键和值
		//key := iter.Key()
		value := iter.Value()
		// 反序列化区块数据
		var tmpBlock Block.Block
		err = tmpBlock.UnmarshalJSON(value)
		if err != nil {
			fmt.Println("err: ", err)
		}
		// 将区块添加到区块链
		tmpBlockchain.chain = append(tmpBlockchain.chain, &tmpBlock)
	}
	iter.Release()

	// 检查遍历过程中是否有错误
	if err := iter.Error(); err != nil {
		return nil, err
	}

	return tmpBlockchain, nil
}
