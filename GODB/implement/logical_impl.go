package implement

import (
	_ "../logical"
	"../physical"
	"../model"
)

type LogicalBase_impl struct{
	storage  *physical.Storage
	tree *model.BinaryTree
	// 节点
	// 值
	// storage
	
}

/*
	commit()
	_refresh_tree_ref() // 更新视图
	get(string) string 
	set(string,string)
	delete()
	follow()
	address()
	len()
	_follow()
*/



func NewLogicalBaseImpl(store *physical.Storage) *LogicalBase_impl{
	return &LogicalBase_impl{
		storage:store,
		tree: model.NewB2Tree(store),
	}
}

func (lb *LogicalBase_impl) get(key string) string{
	/*
	if !storage.Locked{
		lb._refresh_tree_ref()
	}*/
	return lb.tree.Get( lb.tree.Get_root_node(),key)
}

func (lb LogicalBase_impl) set(key string, val string){
	tmp:= lb.tree.Get_root_node()
	lb.tree.Set_new_root( *lb.tree.Insert( &tmp ,key, *model.NewValueRef(-1,val)) )
}

func (lb LogicalBase_impl) commit(){
	lb.tree.Commit()
	//lb.storage.commit_root_address(db.model.address())
}

/* 记录一下 commit的调用路径

->
	tree.Commit(){

		err:=t._tree_ref.Store(t.storage)
		t.storage.Commit_root_address( t._tree_ref.Address() )
	}

->
	_tree_ref is TreeNodeRef | 所有的Ref的store 都会 Prepare_to_store

	TreeNodeRef.Store(){
		if !t._referent.IsNull() && t._address == lg.INIT_ADDR {
			t.Prepare_to_store(storage) // 1. 递归

			t._address = storage.Write(t.Ref2bytes(t._referent)) // 2. 处理本节点
		}
	}

-> 1.
	TreeNodeRef.Prepare_to_store(){
		if !t._referent.IsNull() {
			t._referent.store_refs(st)
		}
	}
->
	_referent is Node

	_referent.store_refs(st){
		left_ref.store(st)
		right_ref.store(st)
		value_ref.store(st)
	}

-> 故可发现是递归下去的


*/


func (lb LogicalBase_impl) _refresh_tree_ref(){}

/*
func (lb LogicalBase_impl)Address() int64 {
	return lb.tree.Address()
}
*/