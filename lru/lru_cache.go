package lru

/*
	1. 实现一个简易的LRU（最近最少使用）缓存。根据数据的历史访问记录来进行淘汰数据，“如果数据最近被访问过，那么将来被访问的几率也更高”。
	2. 使用双向链表，链表节点中包含cache的key和value，此外还有前后指针
	3. 为了加快判断命中效率，使用map缓存key:*node
	4. 最主要两个方法为put和get：
		4.1 put命中将节点移到链表头；
		4.2 put未命中，根据链表节点数是否超过限制来判断是否需要删除链表尾节点，然后在链表头插入节点；
		4.3 get命中，同4.1;
		4.4 get未命中，返回错误。
*/

type Node struct {
	Key, Value int
	Pre, Next  *Node
}

type LRUCache struct {
	Cap        int
	Len        int
	CacheMap   map[int]*Node
	head, tail *Node
}

func InitLRU(cap int) *LRUCache {
	l := &LRUCache{
		Cap:      cap,
		Len:      0,
		CacheMap: map[int]*Node{},
		head:     &Node{},
		tail:     &Node{},
	}
	l.head.Next = l.tail
	l.tail.Pre = l.head
	return l
}

// 根据key获取 value 以及并更新链表
func (l *LRUCache) Get(key int) int {
	if v, ok := l.CacheMap[key]; ok {
		l.moveToHead(v)
		return v.Value
	}
	return -1
}

func (l *LRUCache) Put(key, value int) {
	if v, ok := l.CacheMap[key]; ok {
		v.Value = value
		l.moveToHead(v)
		return
	}
	//不存在
	node := &Node{key, value, nil, nil}
	for l.Len >= l.Cap {
		//删除尾节点
		l.deleteTail()
	}
	l.insertToHead(node)
}

func (l *LRUCache) deleteTail() {
	temp := l.tail.Pre
	temp.Pre.Next = l.tail
	l.tail.Pre = temp.Pre
	temp.Pre = nil
	temp.Next = nil
	l.Len--
	delete(l.CacheMap, l.tail.Key)
}

// 将访问到的节点更新到链表头
func (l *LRUCache) moveToHead(node *Node) {
	node.Pre.Next = node.Next
	node.Next.Pre = node.Pre
	temp := l.head.Next
	l.head.Next = node
	node.Pre = l.head
	node.Next = temp
	temp.Pre = node
}

// 将创建的节点插入链表头
func (l *LRUCache) insertToHead(node *Node) {
	temp := l.head.Next
	l.head.Next = node
	node.Pre = l.head
	node.Next = temp
	temp.Pre = node
	l.Len++
	l.CacheMap[node.Key] = node
}
