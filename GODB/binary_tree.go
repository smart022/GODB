package GODB

import (
	//"fmt"
	"bytes"
	"encoding/gob"
)

// 设计思想
// 除了Ref类变量 其他的结构都是指针型的
// 所有的New 都是生成 指针的

//////////////// NodeRef //////////////
// extends ValueRef
type BinaryNodeRef struct {
	_referent *BinaryNode //
	_address  int64
}

// New
func NewBNodeRef(addr int64, ref *BinaryNode) *BinaryNodeRef {
	ret := new(BinaryNodeRef)
	ret._referent = ref // ref 具体形式是不一样的
	ret._address = addr // 但都可有addr得来，即反序列化
	return ret
}

var NullRef BinaryNodeRef = BinaryNodeRef{nil, init_addr}

func (t *BinaryNodeRef) IsNull() bool {
	return t._address == init_addr
}

func (t *BinaryNodeRef) address() int64 {
	return t._address
}

func (t *BinaryNodeRef) prepare_to_store() {
	// to do
}

func (t *BinaryNodeRef) real(storage *Storage) interface{} { //*BinaryNode
	//init_ref==nil
	if t._referent == nil && t._address != init_addr {
		t._referent = (t.bytes2ref(storage.read(t._address))).(*BinaryNode)
	}
	return t._referent
}

func (t *BinaryNodeRef) store(storage *Storage) {
	// init_ref==nil
	if t._referent != nil && t._address == init_addr {
		t.prepare_to_store()
		t._address = storage.write(t.ref2bytes(t._referent))
	}
}

func (t *BinaryNodeRef) bytes2ref(cot []byte) interface{} {
	var ret BinaryNodeStore
	decoder := gob.NewDecoder(bytes.NewBuffer(cot))
	err := decoder.Decode(&ret)
	check(err)

	return NewNode(
		NewBNodeRef(ret.Left, nil),
		NewBNodeRef(ret.Right, nil),
		NewTrueValRef("", ret.Value),
		ret.Key,
		ret.Length)
}

func (t *BinaryNodeRef) ref2bytes(node interface{}) []byte {

	cur := node.(*BinaryNode)
	//
	sto := &BinaryNodeStore{
		cur.left_ref.address(),
		cur.right_ref.address(),
		cur.value_ref.address(),
		cur.key,
		cur.length}

	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err := encoder.Encode(sto)
	check(err)

	return buf.Bytes()
}

/////////// node //////////////

type BinaryNode struct {
	left_ref  ValueRef
	right_ref ValueRef
	value_ref ValueRef
	key       string
	length    int
}

// use for store seraliza only
type BinaryNodeStore struct {
	Left   int64
	Right  int64
	Value  int64
	Key    string
	Length int
}

//
func NewNode(l, r, v ValueRef, k string, len int) *BinaryNode {
	return &BinaryNode{l, r, v, k, len}
}

// legacy
func FromNode(node *BinaryNode, new_left, new_right BinaryNodeRef, new_val ValueRef) *BinaryNode {
	//l := new_left.(BinaryNodeRef)
	if !new_left.IsNull() {
		node.left_ref = &new_left
		return node
	}
	//r := new_right.(BinaryNodeRef)
	if !new_right.IsNull() {
		node.right_ref = &new_right
		return node
	}

	//rval := new_val.(TrueValRef)
	node.value_ref = new_val
	return node

}

/////////  Treebase ///////////

// extends LogicalBase
type BinaryTree struct {
	storage   *Storage
	_tree_ref BinaryNodeRef
}

func (t *BinaryTree) _refresh_tree_ref() {
	t._tree_ref = *NewBNodeRef(t.storage.get_root_address(), nil)
}

func (t *BinaryTree) get(key string) string {
	if !t.storage.locked {
		t._refresh_tree_ref()
	}
	return t._get((t._follow(&t._tree_ref)).(*BinaryNode), key)
}

func (t *BinaryTree) _get(node *BinaryNode, key string) string {
	for node != nil {
		if key < node.key {
			node = t._follow(node.left_ref).(*BinaryNode)
		} else if key > node.key {
			node = t._follow(node.right_ref).(*BinaryNode)
		} else {
			return t._follow(node.value_ref).(string)
		}
		//panic()
	}
	return ""
}

func (t *BinaryTree) set(key string, val string) {
	if t.storage.lock() {
		t._refresh_tree_ref()
	}
	t._tree_ref = *t._set((t._follow(&t._tree_ref)).(*BinaryNode),
		key,
		NewTrueValRef(val, init_addr))
}

func (t *BinaryTree) _set(root *BinaryNode, key string, val ValueRef) *BinaryNodeRef {
	var new_node *BinaryNode
	if root == nil {
		new_node = NewNode(&NullRef,
			&NullRef,
			val,
			key,
			1)
	} else if key < root.key {
		new_node = FromNode(
			root,
			t._set(t._follow(root.left_ref), key, val),
			&NullRef,
			&NullRef)
	} else if key > root.key {
		new_node = FromNode(
			root,
			&NullRef,
			t._set(t._follow(root.right_ref), key, val),
			&NullRef)

	} else {
		new_node = FromNode(
			root,
			&NullRef,
			&NullRef,
			val)
	}
	return NewBNodeRef(init_addr, new_node)
}

//@ class method
func (t *BinaryTree) _follow(ref ValueRef) interface{} {
	return ref.real(t.storage)
}

func NewB2Tree(storage *Storage) *BinaryTree {
	tree := new(BinaryTree)
	tree.storage = storage
	tree._tree_ref = *NewBNodeRef(storage.get_root_address())
	return tree
}
