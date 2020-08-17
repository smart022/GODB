package implement

import (
	"../logical"
	"../physical"
	"../model"
)

type LogicalBase_impl struct{
	storage  physical.Storage
	tree model.BinaryTree
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
	len()
	_follow()
*/



func NewLogicalBaseImpl(store physical.Storage) *LogicalBase_impl{
	return &LogicalBase_impl{
		storage:store
		tree: 
	}
}

func (lb LogicalBase_impl) get(string key) {
	if !storage.Locked{
		lb._refresh_tree_ref()
	}
	return 
}

func (lb LogicalBase_impl) set(){

}

func (lb LogicalBase_impl) commit(){
	
}
func (lb LogicalBase_impl) _refresh_tree_ref(){}

