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
	//tmp,_:= btree._follow(&btree._tree_ref) // failed!
	//tnode:=tmp.(BinaryNode)
	btree._tree_ref = *btree.Insert(  &NULLNODE, "0", ValueRef{"gigikj",-1})

	t.Log("Insert Tree Succeed!")

}