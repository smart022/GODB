package model
/*
这个 valueref的实现 原来放在 imple下面，但是btree又要用，结果循环依赖了，所以迁移到model下面，感觉反而更合理
*/
import (
	"../logical"
	"../physical"
	"errors"
)

// 实现接口都是基于 指针的，所以要这样 & ...
var _ logical.Ref = &ValueRef{}

// string 实现的值 ref
// extends Ref
type ValueRef struct {
	_referent string // 实际是string ref 包裹
	_address  int64
}

// New
func NewValueRef(addr int64,ref string) *ValueRef {
	return &ValueRef{ref, addr}
}

func (t *ValueRef) Address() int64 {
	return t._address
}


func (t *ValueRef) Prepare_to_store(st *physical.Storage) {
	// do nothing
}


func (t *ValueRef) IsNull() bool {
	return true
}

// string not match nil, so use ""
func (t *ValueRef) Get(st *physical.Storage) (interface{},error){
	var err error
	var ok bool = false
	if t._referent=="" && t._address!= logical.INIT_ADDR{
		t._referent,ok = t.Bytes2ref( st.Read(t._address) ).(string)
	}
	if !ok{
		err = errors.New("[valueref] get type convert failed!")
	}
	return t._referent,err
}

func (t *ValueRef) Store(st *physical.Storage) error{
	var err error
	if t._referent!="" && t._address== logical.INIT_ADDR{
		t.Prepare_to_store(st)
		t._address  = st.Write( t.Ref2bytes(t._referent) )
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
func (t *ValueRef) Bytes2ref(cot []byte) interface{} { // string
	str := string(cot[:])
	return str
}

func (t *ValueRef) Ref2bytes(str interface{}) []byte {
	var data []byte
	if strr,ok:=str.(string);ok{
		data = []byte(strr)
	} else{
		panic("Invalid ref type for ValueRef")
	}
	return data
}

func (t *ValueRef) Length() (int,error) {
	if t._referent==""{
		return 0,errors.New("unloaded ref")
	}

	return int(len(t._referent)), nil
}

/*
	Length() (int64,error) // _referent 的长度
	Prepare_to_store() // 查出 _address 定位
	Address() int64

	// 基础操作
	Get(st *physical.Storage) (interface{},error)
	Store(st *physical.Storage) error
	Bytes2ref([]byte) interface{}
	Ref2bytes(interface{}) []byte // 这个ref2bytes 我才发现存的都是地址
*/