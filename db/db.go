package db

import (
	"github.com/boltdb/bolt"
	"github.com/viviviviviid/go-coin/utils"
)

// bolt는 key와 value만 존재
// bolt는 SQL의 Table과 비슷한 bucket을 갖는다

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil) // Bolt DB 시작, 이름도 생성
		db = dbPointer
		utils.HandleErr(err)
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(dataBucket)) // bucket 생성
			utils.HandleErr(err)
			_, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
			utils.HandleErr(err)
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func Close() { // DB를 계속 열어두면 data 손상이 날 수 있으므로, 꼭 닫아줘야함
	DB().Close()
}

func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data) // db에 저장 key: value => hash: data
		return err
	})
	utils.HandleErr(err)
}

func SaveBlockchain(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data) // db에 저장 key: value => "blockchain": data
		return err
	})
	utils.HandleErr(err)
}

func Checkpoint() []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket)) // dataBucket을 이름으로 하는 버킷을 가져옵니다.
		data = bucket.Get([]byte(checkpoint))   // checkpoint키에 해당하는 값을 가져와서 data 변수에 저장합니다
		return nil                              // 딱히 error을 생성하는 내용이 없기 때문
	})
	return data
}

func Block(hash string) []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}
