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

	Length() (int,error) // _referent 的长度
	Prepare_to_store(st *physical.Storage) // 查出 _address 定位
	Address() int64

	// 基础操作
	Get(st *physical.Storage) (interface{},error)
	Store(st *physical.Storage) error
	Bytes2ref([]byte) interface{}
	Ref2bytes(interface{}) []byte // 这个ref2bytes 我才发现存的都是地址
}

const (
	//INIT_REF  = ""
	INIT_ADDR = -1
)