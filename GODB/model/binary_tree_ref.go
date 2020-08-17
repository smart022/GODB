package model

import (
	//"fmt"
	"bytes"
	"encoding/gob"
	"../physical"
	lg "../logical"
)

// 这个tree_ref 是 node_ref 也应该归于 impl 的，毕竟impl是logical的实现层，
// 但这个归属还是忙模糊的， 因为model/tree也是logical的某种实现。。

//////////////// NodeRef //////////////
// extends Ref
type BinaryNodeRef struct {
	_referent BinaryNode //
	_address  int64
}

var _ lg.Ref = NewBNodeRef(0,nil)

// New
func NewBNodeRef(addr int64, ref BinaryNode) BinaryNodeRef {
	
	return BinaryNodeRef{
		_referent:ref
		_address: addr
	}
}

var NullRef BinaryNodeRef = BinaryNodeRef{nil, init_addr}

func (t *BinaryNodeRef) IsNull() bool {
	return t._address == lg.INIT_ADDR
}

func (t *BinaryNodeRef) address() int64 {
	return t._address
}

// 注意  prepare_to_store 和 ref . store 的交互调用
// 保证了树的递归存储
func (t *BinaryNodeRef) prepare_to_store(st *physical.Storage) {
	if t._referent!=nil{
		t._referent.store_refs(st)
	}
}

func (t *BinaryNodeRef) get(storage *physical.Storage) (interface{},error) { //*BinaryNode
	
	var err error
	if t._referent == nil && t._address != lg.INIT_ADDR {
		t._referent = (t.bytes2ref(storage.read(t._address))).(BinaryNode)
	}

	return t._referent,err
}

func (t *BinaryNodeRef) store(storage *physical.Storage) error {
	
	var err error
	if t._referent != nil && t._address == lg.INIT_ADDR {
		t.prepare_to_store(storage)
		t._address = storage.write(t.ref2bytes(t._referent))
	}

	return err
}

func (t *BinaryNodeRef) bytes2ref(cot []byte) interface{} {
	var ret BinaryNodeStore
	decoder := gob.NewDecoder(bytes.NewBuffer(cot))
	err := decoder.Decode(&ret)

	panic_check(err,"[BinaryNodeRef] bytes2ref failed! Invalid bytes")

	return NewNode(
		NewBNodeRef(ret.Left, nil),
		NewBNodeRef(ret.Right, nil),
		NewTrueValRef(ret.Value, ""),
		ret.Key,
		ret.Length)
}


// 这个
func (t *BinaryNodeRef) ref2bytes(node interface{}) []byte {

	cur := node.(BinaryNode)
	//
	sto := BinaryNodeStore{
		cur.left_ref.address(),
		cur.right_ref.address(),
		cur.value_ref.address(),
		cur.key,
		cur.length}

	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err := encoder.Encode(sto)

	panic_check(err,"[BinaryNodeRef] ref2bytes failed!")

	return buf.Bytes()
}

func panic_check(err error, info string){
	if err!=nil{
		panic(info)
	}
}