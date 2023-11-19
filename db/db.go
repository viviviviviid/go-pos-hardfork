// db 패키지는 db 조회, 수정 등과 같은 함수를 제공합니다.
package db

import (
	"fmt"
	"os"

	"github.com/viviviviviid/go-coin/utils"
	bolt "go.etcd.io/bbolt"
)

// bolt는 key와 value만 존재
// bolt는 SQL의 Table과 비슷한 bucket을 갖는다
const (
	dbName       = "blockchain"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

var db *bolt.DB

type DB struct{}

func (DB) FindBlock(hash string) []byte {
	return findBlock(hash)
}
func (DB) LoadChain() []byte {
	return loadChain()
}
func (DB) SaveBlock(hash string, data []byte) {
	saveBlock(hash, data)
}
func (DB) SaveChain(data []byte) {
	saveChain(data)
}
func (DB) DeleteAllBlocks() {
	emptyBlocks()
}

// 노드 포트 번호를 이용하여 DB를 탐색 (Ex. blockchain_4000.db)
func getDbName() string {
	port := os.Args[2][6:]
	return fmt.Sprintf("./node_dbs/%s_%s.db", dbName, port)
}

// 노드 실행시 DB 유무 확인 후 생성 또는 불러오기
func InitDB() {
	if db == nil {
		dbPointer, err := bolt.Open(getDbName(), 0600, nil) // Bolt DB 시작, 이름도 생성
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
}

// DB의 data 손상을 막기 위해 DB 닫기
func Close() {
	db.Close()
}

// DB내에 블록 데이터 저장
func saveBlock(hash string, data []byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data) // key: value = hash: data
		return err
	})
	utils.HandleErr(err)
}

// DB내에 체인 데이터 저장
func saveChain(data []byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data) // key: value = blockchain: data
		return err
	})
	utils.HandleErr(err)
}

// DB내의 체인 정보 체크포인트까지 불러오기
func loadChain() []byte {
	var data []byte
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

// DB내에서 특정 블록 검색
func findBlock(hash string) []byte {
	var data []byte
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}

// 모든 블록 제거 (노드간 블록높이 비교 후 대체 할때 사용)
func emptyBlocks() {
	db.Update(func(tx *bolt.Tx) error {
		utils.HandleErr(tx.DeleteBucket([]byte(blocksBucket)))
		_, err := tx.CreateBucket([]byte(blocksBucket))
		utils.HandleErr(err)
		return nil
	})
}
