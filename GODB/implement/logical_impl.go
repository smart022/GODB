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
func (lb LogicalBase_impl) _refresh_tree_ref(){}

/*
func (lb LogicalBase_impl)Address() int64 {
	return lb.tree.Address()
}
*/