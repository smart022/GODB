package GODB

import (
	"log"
	"os"
	"./physical"
	"./model"
	"./logical"
	"./implement"
)

type GODB struct {
	storage *physical.Storage     ``
	model   *logical.LogicalBase ``
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

func (db *GODB) commit()  {
	if !db.storage.assert_not_closed() {
		return error.New("file closed")
	}
	db.storage.commit_root_address(db.model.address())
	return nil
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

	// 基础实现源 变成logicalbase_impl
	db.model = NewLogicalBaseImpl( db.storage )
	return db
}
