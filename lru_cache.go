package local_cache

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
	"time"
)

type Node struct {
	pre  *Node
	next *Node
	key  string
	val  item
}

func NewNode(key string, val item) *Node {
	return &Node{
		pre:  nil,
		next: nil,
		key:  key,
		val:  val,
	}
}

type NodeList struct {
	head *Node //最近使用
	end  *Node //最久未使用
}

func NewNodeList() *NodeList {
	head := NewNode("head", item{})
	end := NewNode("end", item{})

	head.next = end
	end.pre = head
	return &NodeList{
		head: head,
		end:  end,
	}
}

func (n *NodeList) front() *Node {
	if n.isEmpty() {
		return nil
	}

	return n.head.next
}

func (n *NodeList) back() *Node {
	if n.isEmpty() {
		return nil
	}

	return n.end.pre
}

func (n *NodeList) isEmpty() bool {
	return n.head.next == n.end
}

func (n *NodeList) pushFront(node *Node) {
	node.next = n.head.next
	node.next.pre = node
	n.head.next = node
	node.pre = n.head
}

func (n *NodeList) popBack() {
	node := n.end.pre
	if node == n.head {
		return
	}

	node.pre.next = n.end
	n.end.pre = node.pre
}

func (n *NodeList) moveToFront(node *Node) {
	node.next.pre = node.pre
	node.pre.next = node.next

	node.pre = n.head
	node.next = n.head.next
	node.next.pre = node
	n.head.next = node
}

func (n *NodeList) moveToBack(node *Node) {
	node.next.pre = node.pre
	node.pre.next = node.next

	node.next = n.end
	node.pre = n.end.pre
	node.pre.next = node
	n.end.pre = node
}

func (n *NodeList) delete(node *Node) {
	node.next.pre = node.pre
	node.pre.next = node.next
}

type LRUCache struct {
	valueMap  map[string]*Node
	options   *Options
	rwMutex   sync.RWMutex
	sf        singleflight.Group
	list      *NodeList
	cleanFLag bool
}

func newLRUCache(options *Options) *LRUCache {
	return &LRUCache{
		valueMap: make(map[string]*Node),
		options:  options,
		list:     NewNodeList(),
	}
}

func (l *LRUCache) Get(key string) (interface{}, error) {
	node, err, _ := l.sf.Do(key, func() (interface{}, error) {
		l.rwMutex.RLock()
		defer l.rwMutex.RUnlock()
		if val, exist := l.valueMap[key]; exist {
			return val, nil
		}

		return nil, ErrKeyNotExist
	})

	if err != nil {
		return nil, err
	}

	if _, ok := node.(*Node); !ok {
		return nil, ErrKeyValue
	}
	res := node.(*Node)

	if res.val.isExpired() {
		return nil, ErrKeyExpired
	}

	l.list.moveToFront(res)

	return res.val, nil
}

func (l *LRUCache) Set(key string, val interface{}, opts ...OptionsFn) {
	o := getItemDefaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	node := NewNode(key, item{
		value:      val,
		expiration: time.Now().Add(o.expiration),
	})
	if int64(len(l.valueMap)) > l.options.maxSum && !l.isCleanIng(){
		l.Clean()
	}
	fmt.Println(len(l.valueMap))

	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()
	l.list.pushFront(node)
	l.valueMap[key] = node
}

func (l *LRUCache) Delete(key string) {
	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()
	if node, ok := l.valueMap[key]; ok {
		l.list.delete(node)
		delete(l.valueMap, key)
	}
}

func (l *LRUCache) Clean() {
	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()

	l.cleanFLag = true

	remain := l.getRemain()
	fmt.Println("开始清理")

	for remain > 0 {
		node := l.list.end.pre
		l.list.delete(node)
		delete(l.valueMap, node.key)
		remain--
	}
	l.cleanFLag = false
}

func (l *LRUCache) getRemain() int {
	return len(l.valueMap) / 2
}

func (l *LRUCache) isCleanIng() bool {
	return l.cleanFLag
}

func (l *LRUCache) Close() {
}
