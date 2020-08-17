package model

import (
	//"fmt"
	"bytes"
	"encoding/gob"
	"../physical"
	lg "../logical"
)

/////////// node //////////////

type BinaryNode struct {
	left_ref  lg.Ref
	right_ref lg.Ref
	value_ref lg.Ref
	key       string
	length    int // 这个是 只树的总节点数
}

// use for store seraliza only
type BinaryNodeStore struct {
	// 都是地址
	Left   int64
	Right  int64
	Value  int64

	Key    string
	Length int
}

const (
	NULLREF lg.Ref = nil
)

//
func NewNode(l, r, v lg.Ref, k string, len int) *BinaryNode {
	return &BinaryNode{l, r, v, k, len}
}

// FromNode 的目的是为了实现insert 不改变原树，返回全新的树，但共用原来的节点数据
func FromNode(node , new_left, new_right BinaryNodeRef, new_val ValueRef) BinaryNode {
	

	// 只改了一个节点的值，所以长度不变
	if new_val!=nil {
		return BinaryNode{
			left_ref : node.left_ref,
			right_ref : node.right_ref,
			value_ref: new_val,
			key: node.key,
			length: node.length
		}
	}

	ori_len := node.length

	// 因为insert的逻辑中，只会插左右中的一个，所以进入一个后就可以return
	if !new_left.IsNull() {

		
		return BinaryNode{
			left_ref : new_left,
			right_ref : node.right_ref,
			value_ref: node.value_ref,
			key: node.key,
			length: ori_len + ( new_left.length -node.left_ref.length) // 注意这个
			// 有2点
			// 1 是语义: 原根len + 新子树len -老子树len
			// 2 是计算的时候发现： len有node的 也有noderef的 这里混用的，大概是 noderef和 node的len的意义是一样的
		}
	}
	//r := new_right.(BinaryNodeRef)
	if !new_right.IsNull() {
		return BinaryNode{
			left_ref : node.left_ref,
			right_ref : new_right,
			value_ref: node.value_ref,
			key: node.key,
			length: ori_len + ( new_right.length -node.right_ref.length)
		}
	}

	return nil
}

/////////  Treebase ///////////

// extends LogicalBase
// 树在源码中是 LogicalBase 的子类， 但在golang中不把其实现为之类， 只提供逻辑相同的接口，在logical_impl中被调用
type BinaryTree struct {
	storage   *Storage
	_tree_ref BinaryNodeRef // 本质上ref是一个节点的包裹，所以root就是个ref
}

func (t *BinaryTree) _refresh_tree_ref() {
	t._tree_ref = *NewBNodeRef(t.storage.get_root_address(), nil)
}

/* impl中实现
func (t *BinaryTree) get(key string) string {
	if !t.storage.locked {
		t._refresh_tree_ref()
	}
	return t._get((t._follow(&t._tree_ref)).(*BinaryNode), key)
}
*/

func (t *BinaryTree) _get(node BinaryNode, key string) string {
	var err error
	for node != nil && err==nil {
		if key < node.key {
			node,err = t._follow(node.left_ref).(BinaryNode)
		} else if key > node.key {
			node,err = t._follow(node.right_ref).(BinaryNode)
		} else {
			ret,nerr := t._follow(node.value_ref).(string)
			
			if nerr!=nil{
				return ""
			}

			return ret
		}
		
	}
	return ""
}

/* impl中实现
func (t *BinaryTree) set(key string, val string) {
	if t.storage.lock() {
		t._refresh_tree_ref()
	}
	t._tree_ref = *t._set((t._follow(&t._tree_ref)).(*BinaryNode),
		key,
		NewTrueValRef(val, init_addr))
}
*/

// _insert return 是 BinaryNodeRef
// 但中 root 是 node
func (t *BinaryTree) _insert(root *BinaryNode, key string, val ValueRef) BinaryNodeRef {
	var new_node *BinaryNode

	if root == nil {
		new_node = NewNode(NULLREF,
			NULLREF,
			val,
			key,
			1)
	} else if key < root.key {
		new_node = FromNode(
			root,
			t._insert(t._follow(root.left_ref), key, val),
			&NullRef,
			&NullRef)
	} else if key > root.key {
		new_node = FromNode(
			root,
			&NullRef,
			t._insert(t._follow(root.right_ref), key, val),
			&NullRef)

	} else {
		new_node = FromNode(
			root,
			&NullRef,
			&NullRef,
			val)
	}
	// 返回这个插入的 nodeRef
	return NewBNodeRef(lg.INIT_ADDR, new_node)
}

//@ class method
// 这个_follow给人一种 后续节点的感觉，但其实不是，这个其实就是ref解包，获得其_referent实体
func (t *BinaryTree) _follow(ref lg.Ref) (interface{},error) {
	return ref.get(t.storage)
}

func NewB2Tree(storage *Storage) *BinaryTree {
	tree := new(BinaryTree)
	tree.storage = storage
	tree._tree_ref = *NewBNodeRef(storage.get_root_address())
	return tree
}
