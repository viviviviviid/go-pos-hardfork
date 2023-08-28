package db

import (
	"github.com/boltdb/bolt"
	"github.com/viviviviviid/go-coin/utils"
)

// bolt는 key와 value만 존재
// bolt는 SQL의 Table과 비슷한 bucket을 갖는다

const (
	dbName       = "blockchaing.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil) // Bolt DB 시작, 이름도 생성
		db = dbPointer
		utils.HandleErr(err)
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket)) // bucket 생성
			utils.HandleErr(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			utils.HandleErr(err)
			return err
		})
		utils.HandleErr(err)
	}
	return db
}
