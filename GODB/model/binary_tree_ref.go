package model

import (
	//"fmt"
	"bytes"
	"encoding/gob"
	"errors"
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

var _ lg.Ref = NewBNodeRef(0,BinaryNode{})

// New
func NewBNodeRef(addr int64, referent BinaryNode) *BinaryNodeRef {
	
	return &BinaryNodeRef{
		_referent:referent,
		_address: addr,
	}
}

//var NullRef BinaryNodeRef = BinaryNodeRef{nil, init_addr}

func (t *BinaryNodeRef) Address() int64 {
	return t._address
}

func (t *BinaryNodeRef) IsNull() bool {
	// 要考虑一下
	return t._referent.IsNull()
	//return t._address == lg.INIT_ADDR
}

func (t *BinaryNodeRef) Length() (int,error) {
	if t._referent.IsNull() {
		return 0,errors.New("unloaded ref")
	}

	return int(t._referent.length), nil
}

// 注意  prepare_to_store 和 ref . store 的交互调用
// 保证了树的递归存储
func (t *BinaryNodeRef) Prepare_to_store(st *physical.Storage) {
	if !t._referent.IsNull() {
		t._referent.store_refs(st)
	}
}

func (t *BinaryNodeRef) Get(storage *physical.Storage) (interface{},error) { //*BinaryNode
	
	var err error
	if t._referent.IsNull()  && t._address != lg.INIT_ADDR {
		t._referent = (t.Bytes2ref(storage.Read(t._address))).(BinaryNode)
	}

	return t._referent,err
}

func (t *BinaryNodeRef) Store(storage *physical.Storage) error {
	
	var err error
	if !t._referent.IsNull() && t._address == lg.INIT_ADDR {
		t.Prepare_to_store(storage)
		t._address = storage.Write(t.Ref2bytes(t._referent))
	}

	return err
}

func (t *BinaryNodeRef) Bytes2ref(cot []byte) interface{} {
	var ret BinaryNodeStore
	decoder := gob.NewDecoder(bytes.NewBuffer(cot))
	err := decoder.Decode(&ret)

	panic_check(err,"[BinaryNodeRef] bytes2ref failed! Invalid bytes")

	return NewNode(
		NewBNodeRef(ret.Left, NULLNODE),
		NewBNodeRef(ret.Right, NULLNODE),
		NewValueRef(ret.Value, ""),
		ret.Key,
		ret.Length)
}


// 这个
func (t *BinaryNodeRef) Ref2bytes(node interface{}) []byte {

	cur := node.(BinaryNode)
	//
	sto := BinaryNodeStore{
		cur.left_ref.Address(),
		cur.right_ref.Address(),
		cur.value_ref.Address(),
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