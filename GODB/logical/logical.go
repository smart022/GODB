package logical

// 树结构 模型基础操作的抽象， 是DB内成员，供给原始接口


/* 源python 中

	db.get(key) -> db.tree.get(key)
	而 db.tree 是 logicalbase的子类
	logicalbase有 get 里面带简单锁判断逻辑 然后 调用 tree._get() 这个子类的函数？？来返回
	问题就是 父类默认有子类这个_get接口 ， 这个类型设计是不是有点迷惑。

	在go的实现版本中不能照单全收，要改进一下
	我现在要拆解一下
	db的所有操作 都走 logicalbase_impl 这个实现类， 而实现类里带 tree实例

	所以python源码的 logicalbase 实际上 本项目的 logicalbase_impl


	// ref 接口概念，( 了解后， 发现ref是个 实体的包裹， 使用了lazy init， 因为有地址就行了，需要实体的时候取出来就好了
	有成员 _referent, _address

	_referent 应该是一个实体元素，值，比如string,也可以说复合的 struct{}
	_address 是这个实体元素在文件中的地址

	实现了 address, referent2string, length
	get,store 等操作 就是数值变二进制，在文件中定位

	// value_ref 和 node_ref
	// 
	node {
		key
		val_ref 

		left_ref
		right_ref
	}

	// prepare_to_store: _address 也可以为空，延迟查找


	// 而在btree的get中，get操作是获得这个ref的实体 _referent，
	比较的方法是 node = _follow(node.left_ref)
	_follow (ref) = return ref.get(self._storage)

	// valueRef 的get ：目的是延迟获取
		
	def get(self, storage):
        if self._referent is None and self._address:
            self._referent = self.string_to_referent(storage.read(self._address))
		return self._referent
		

	// _follow 是 logicalbase的成员函数, 
*/
// 因为go的接口抽象都是纯虚, 但参考的源码是python写的，所以其内部是有部分逻辑实现的所以需要加一层来实现部分逻辑
// 在impl

////////////////////////////
// model //

type LogicalBase interface {
	//
	// node_ref_class = None
	// value_ref_class = ValueRef

	// methods
	Commit()
	Refresh_tree_ref() // 更新视图
	Get(string) string 
	Set(string,string)
	delete()
	//follow()
	Len() int
	Address() int64
	_follow()
}
