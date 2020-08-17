package logical

import "../physical"

// 解释见 logical.go
// 节点值类型 抽象

/*
	成员变量
	_referent interface{}
	_address int64
*/
type Ref interface {

	length() (int64,error) // _referent 的长度
	prepare_to_store() // 查出 _address 定位
	address() int64

	// 基础操作
	get(st *physical.Storage) (interface{},error)
	store(st *physical.Storage) error
	bytes2ref([]byte) interface{}
	ref2bytes(interface{}) []byte // 这个ref2bytes 我才发现存的都是地址
}

const (
	//INIT_REF  = ""
	INIT_ADDR = -1
)