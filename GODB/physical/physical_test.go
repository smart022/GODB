package physical

import (
	"os"
	"testing"
)

func TestFunc(t *testing.T) {
	//t.Skip()
	//t.Errorf(os.Getwd())
	addr := "test1.txt"
	f, err := os.OpenFile(addr, os.O_CREATE|os.O_RDWR, 0755)
	if err!=nil{
		panic(err.Error())
	}

	st := NewStorage(f)
	if st == nil {
		t.Errorf("file open failed!")
	}

	//t.Log("End pos: ", st._seek_end())

	bt := make([]byte, 5)
	for i := 0; i < 5; i++ {
		bt[i] = 0x12
	}
	pos := st.write(bt)
	t.Log("write pos: ", pos)
	t.Log("write res ", bt)

	st._f.Sync()

	//st._f.Seek(0, 0)
	tmp := st.read(0)
	t.Log("get res ", tmp)

	t.Log("pass!")
}

func TestInteger(t *testing.T) {
	//t.Skip()
	addr := "test.txt"
	f, err := os.OpenFile(addr, os.O_CREATE|os.O_RDWR, 0755)
	if err!=nil{
		panic(err.Error())
	}
	defer f.Close()	
	st := NewStorage(f)
	if st == nil {
		t.Errorf("file open failed!")
	}
	var num int64

	num = 1432

	err = st.write_integer(num)
	if err != nil {

		t.Errorf("%s %s", err.Error(),"write_integer failed!")
	}

	st._f.Sync()


	// reopen

	st._f.Seek(0, 0)
	ret := st.read_integer()
	if ret != num {
		t.Errorf("read_integer failed!")
	}
	t.Log(num, ret)
	t.Log("Integer pass!")
}
