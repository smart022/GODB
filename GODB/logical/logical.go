package logical

// 先虚
type ValueRef interface {
	real(*Storage) interface{} //`referent`
	store(*Storage)
	prepare_to_store()
	address() int64

	//
	bytes2ref([]byte) interface{}
	ref2bytes(interface{}) []byte
}

const (
	init_ref  = ""
	init_addr = -1
)

// 再实
// extends ValueRef
type TrueValRef struct {
	_referent string //
	_address  int64
}

// New
func NewTrueValRef(ref string, addr int64) *TrueValRef {
	ret := TrueValRef{ref, addr}
	//ret._referent = ref // ref 具体形式是不一样的
	//ret._address = addr // 但都可有addr得来，即反序列化
	return &ret
}

func (t *TrueValRef) address() int64 {
	return t._address
}

func (t *TrueValRef) real(storage *Storage) interface{} { // string
	if t._referent == init_ref && t._address != init_addr {
		t._referent = (t.bytes2ref(storage.read(t._address))).(string)
	}
	return t._referent
}

func (t *TrueValRef) store(storage *Storage) {
	if t._referent != init_ref && t._address == init_addr {
		t.prepare_to_store()
		t._address = storage.write(t.ref2bytes(t._referent))
	}
}

func (t *TrueValRef) prepare_to_store() {
	// do nothing
}

// 转化的方法与具体子类类型 不一样
// 本应该是类方法
// 但go似乎不支持类方法
func (t *TrueValRef) bytes2ref(cot []byte) interface{} { // string
	str := string(cot[:])
	return str
}

func (t *TrueValRef) ref2bytes(str interface{}) []byte {
	strr := str.(string)
	data := []byte(strr)
	return data
}

////////////////////////////
// model //

type LogicalBase interface {
	//
	// node_ref_class = None
	// value_ref_class = ValueRef

	// methods
	commit()
	_refresh_tree_ref()
	get(string) string
	set(string,string)
	delete()
	follow()
	len()
	_follow()
}
