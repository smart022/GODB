package model


import (
	"os"
	"testing"
	"../physical"
)



func TestBTree(t *testing.T){
	addr := "test.txt"
	f, err := os.OpenFile(addr, os.O_CREATE|os.O_RDWR, 0755)
	defer f.Close()
	if err!=nil{
		panic(err.Error())
	}

	st := physical.NewStorage(f)
	if st == nil {
		t.Errorf("file open failed!")
	}
	btree:=NewB2Tree(st)

	if btree == nil{
		t.Errorf("new tree failed!")
	}

	t.Log("Build Tree Succeed!")

	// 明明在内存里为啥压查文件
	// 这里要伪实现一个 logi_impl set

	btree._tree_ref = *btree.Insert(  &NULLNODE, "0", ValueRef{"gigikj",-1})
	btree.Commit()

	t.Log("Insert Tree Succeed!")

}



func TestBTreeRef(t *testing.T){

	var (
		valref= ValueRef{"aaaaa",-1}
	)

	var node_core *BinaryNode = NewNode( &NULLNODEREF,
		&NULLNODEREF,
		valref,
		"key",
		1)

	var noderef *BinaryNodeRef = NewBNodeRef(0,*node_core)

	if noderef.Address()!=0{
		t.Errorf("Address failed!")
	}

	if noderef.IsNull(){
		t.Errorf("IsNull failed!")
	}

	tmp_bytes:=noderef.Ref2bytes(node_core)
	t.Log(tmp_bytes)

}