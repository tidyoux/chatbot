package db

import (
	"crypto/sha256"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	dbPath = "db/"
)

var (
	dbInstance *leveldb.DB
)

func init() {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		panic(err)
	}
	dbInstance = db
}

type DB struct {
	namespace []byte
}

func New(namespace string) *DB {
	hash := sha256.Sum256([]byte(namespace))
	return &DB{
		namespace: hash[:],
	}
}

func (db *DB) key(key []byte) []byte {
	return append(db.namespace, key...)
}

func (db *DB) Get(key []byte) ([]byte, error) {
	value, err := dbInstance.Get(db.key(key), nil)
	if err == nil {
		return value, nil
	}

	if err == leveldb.ErrNotFound {
		return nil, nil
	}

	return nil, err
}

func (db *DB) Set(key, value []byte) error {
	return dbInstance.Put(db.key(key), value, nil)
}
