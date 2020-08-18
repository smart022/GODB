package model

import (
	//"fmt"
	//"bytes"
	//"encoding/gob"
	"../physical"
	lg "../logical"
	//impl "../implement"
)

/////////// node //////////////

type BinaryNode struct {
	left_ref  lg.Ref
	right_ref lg.Ref
	value_ref lg.Ref
	key       string
	length    int // 这个是 子树的总节点数
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

var (
	NULLVALREF = ValueRef{"",-1} //NewValueRef(-1,"")
	
	NULLNODE = BinaryNode{length:0}
	NULLNODEREF  = BinaryNodeRef{_referent:NULLNODE, _address:-1}
)

func (node *BinaryNode)store_refs(st *physical.Storage){
	node.left_ref.Store(st)
	node.right_ref.Store(st)
	node.value_ref.Store(st)
}

func (node BinaryNode) IsNull() bool {
	return node.length==0
}

// 把Ref设成 interface{} 防止转得麻烦
// l,r,v 都必须是 * xxxRef 即指针
func NewNode(l, r, v interface{}, k string, len int) *BinaryNode {
	return &BinaryNode{ l.(*BinaryNodeRef), 
		r.(*BinaryNodeRef), 
		v.(*ValueRef), 
		k, len}
}

// FromNode 的目的是为了实现insert 不改变原树，返回全新的树，但共用原来的节点数据
func FromNode(node BinaryNode , new_left, new_right *BinaryNodeRef, new_val *ValueRef) *BinaryNode {
	

	// 只改了一个节点的值，所以长度不变
	if !new_val.IsNull() {
		return &BinaryNode{
			left_ref : node.left_ref,
			right_ref : node.right_ref,
			value_ref: new_val,
			key: node.key,
			length: node.length}
	}

	ori_len := node.length

	// 因为insert的逻辑中，只会插左右中的一个，所以进入一个后就可以return
	if !new_left.IsNull() {

		l_len,_:=new_left.Length()
		ol_len,_:= node.left_ref.Length()
		return &BinaryNode{
			left_ref : new_left,
			right_ref : node.right_ref,
			value_ref: node.value_ref,
			key: node.key,
			length: ori_len + (l_len - ol_len ), // 注意这个
			// 有2点
			// 1 是语义: 原根len + 新子树len -老子树len
			// 2 是计算的时候发现： len有node的 也有noderef的 这里混用的，大概是 noderef和 node的len的意义是一样的
		}
	}
	//r := new_right.(BinaryNodeRef)
	if !new_right.IsNull() {
		r_len,_:=new_right.Length() 
		or_len,_:= node.right_ref.Length()
		return &BinaryNode{
			left_ref : node.left_ref,
			right_ref : new_right,
			value_ref: node.value_ref,
			key: node.key,
			length: ori_len + ( r_len - or_len),
		}
	}

	return nil
}

/////////  Treebase ///////////

// extends LogicalBase
// 树在源码中是 LogicalBase 的子类， 但在golang中不把其实现为之类， 只提供逻辑相同的接口，在logical_impl中被调用
type BinaryTree struct {
	storage   *physical.Storage
	_tree_ref BinaryNodeRef // 本质上ref是一个节点的包裹，所以root就是个ref
}

func (t *BinaryTree) _refresh_tree_ref() {
	t._tree_ref = *NewBNodeRef(t.storage.Get_root_address(), NULLNODE)
}


func (t *BinaryTree) Get(node BinaryNode, key string) string {
	var err error
	var tmp interface{}
	// 这个node.isnull() 设计成实例挂载 而不是指针挂载
	for node.IsNull() && err==nil {
		if key < node.key {
			tmp,err = t._follow(node.left_ref)
			node = tmp.(BinaryNode)
		} else if key > node.key {
			tmp,err = t._follow(node.right_ref)
			node = tmp.(BinaryNode)
		} else {
			tmp,nerr := t._follow(node.value_ref)
			ret := tmp.(string)
			if nerr!=nil{
				return ""
			}

			return ret
		}
		
	}
	return ""
}

// _insert return 是 BinaryNodeRef
// 但中 root 是 node
// 这个返回的是一个全新的子树
// 源代码的 逻辑 写到这里很不清晰 ， 需要重新设计
/*

self._tree_ref = self._insert(
	self._follow(self._tree_ref), key, self.value_ref_class(value))
*/
func (t *BinaryTree) Insert(root *BinaryNode, key string, val ValueRef) *BinaryNodeRef {
	var new_node *BinaryNode

	if root.IsNull() {
		new_node = NewNode( &NULLNODEREF,
			&NULLNODEREF,
			&val,
			key,
			1)
	} else if key < root.key {
		tmp,err:=t._follow(root.left_ref)
		check(err)
		bnode_tmp:=tmp.(BinaryNode)
		
		// FromNode(node BinaryNode , new_left, new_right *BinaryNodeRef, new_val *ValueRef) *BinaryNode
		new_node = FromNode(
			*root,
			t.Insert( &bnode_tmp, key, val),
			&NULLNODEREF,
			&NULLVALREF)
	} else if key > root.key {
		tmp,err:=t._follow(root.right_ref)
		check(err)
		bnode_tmp:=tmp.(BinaryNode)
		
		new_node = FromNode(
			*root,
			&NULLNODEREF,
			t.Insert(&bnode_tmp, key, val),
			&NULLVALREF)

	} else {
		new_node = FromNode(
			*root,
			&NULLNODEREF,
			&NULLNODEREF,
			&val)
	}
	// 返回这个插入的 nodeRef
	// 是不是addr 0
	return NewBNodeRef(0, *new_node)
}

//@ class method
// 这个_follow给人一种 后续节点的感觉，但其实不是，这个其实就是ref解包，获得其_referent实体
func (t *BinaryTree) _follow(ref lg.Ref) (interface{},error) {
	return ref.Get(t.storage)
}


// 源码是 logicalbase的逻辑， 这里 tree和 logi_impl都要实现， 然后实际是 logi_impl调 tree
func (t *BinaryTree) Commit(){
	err:=t._tree_ref.Store(t.storage)
	check(err)
	t.storage.Commit_root_address( t._tree_ref.Address() )
}

// 解决麻烦的取rootnode的问题
func (t *BinaryTree) Get_root_node() BinaryNode{
	if t._tree_ref._referent.IsNull(){
		tmp,_:=t._follow( &t._tree_ref )
		ret:=tmp.(BinaryNode)
		return ret
	}
	return  t._tree_ref._referent
}

func (t *BinaryTree) Set_new_root(root BinaryNodeRef){
	t._tree_ref = root //*NewBNodeRef( 0,root)
}

func NewB2Tree(storage *physical.Storage) *BinaryTree {
	tree := new(BinaryTree)
	tree.storage = storage
	tree._tree_ref = *NewBNodeRef(0,NULLNODE)
	return tree
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}