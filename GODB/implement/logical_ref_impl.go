package implement

import "../logical"


var _ logical.Ref = ValueRef{}

// string 实现的值 ref
// extends Ref
type ValueRef struct {
	_referent string // 实际是string ref 包裹
	_address  int64
}

// New
func NewValueRef(addr int64,ref string) TrueValRef {
	return TrueValRef{ref, addr}
}

func (t *ValueRef) address() int64 {
	return t._address
}


func (t *ValueRef) prepare_to_store() {
	// do nothing
}

// string not match nil, so use ""
func (t *ValueRef) get(st *physical.Storage) (interface{},error){
	var err error
	if t._referent=="" && t._address!= logical.INIT_ADDR{
		t._referent ,err = t.bytes2ref( st.read(t._address) )
	}
	return t._referent,err
}

func (t *ValueRef) store(st *physical.Storage) error{
	var err error
	if t._referent!="" && t._address== logical.INIT_ADDR{
		t.prepare_to_store(st)
		t._address ,err = st.write( t.ref2bytes(t._referent) )
	}
	/*
	else{
		err = Error()
	}
	*/

	return err
}


// 转化的方法与具体子类类型 不一样
// 本应该是类方法
// 但go似乎不支持类方法
func (t *ValueRef) bytes2ref(cot []byte) interface{} { // string
	str := string(cot[:])
	return str
}

func (t *ValueRef) ref2bytes(str interface{}) []byte {
	if strr,ok:=str.(string);ok{
		data := []byte()
	} else{
		panic("Invalid ref type for ValueRef")
	}
	return data
}

func (t *ValueRef)length() (int64,error) {
	if t._referent==""{
		return 0,erros.New("unloaded ref")
	}

	return int64(len(t._referent)), nil
}