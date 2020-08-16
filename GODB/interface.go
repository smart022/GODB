package GODB

import (
	"log"
	"os"
)

type GODB struct {
	storage *Storage     ``
	model   *LogicalBase ``
}

func (db *GODB) get(key string) (string, bool) {
	if !db.storage.assert_not_closed() {
		return nil, false
	}

	return model.get(key)
}

func (db *GODB) set(key string, val string) bool {
	if !db.storage.assert_not_closed() {
		return false
	}
	db.model.set(key,val)
	return true
}

func (db *GODB) del(key string) bool {
	return true
}

func (db *GODB) commit() bool {
	return true
}

func NewGODB(addr string) *GODB {
	db := new(GODB)

	f, err := os.OpenFile(addr, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// defer f.Close() 何时close?
	db.storage = NewStorage(f)
	db.model = NewB2Tree()
	return db
}
