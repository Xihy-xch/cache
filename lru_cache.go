package local_cache

import (
	"fmt"
	"sync"
	"time"
)

type LRUCache struct {
	valueMap  map[string]*Node
	maxSum    int64
	rwMutex   sync.RWMutex
	list      *NodeList
	cleanFLag bool
}

func NewLRUCache(maxSum int64) *LRUCache {
	return &LRUCache{
		valueMap: make(map[string]*Node),
		list:     NewNodeList(),
		maxSum:   maxSum,
	}
}

func (l *LRUCache) Get(key string, opts ...OptionsFn) (interface{}, error) {
	o := getDefaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	val, err := l.doGet(key)
	if err == nil {
		return val, nil
	}

	if o.getter == nil {
		return nil, err
	}

	val, err, _ = sf.Do(key, func() (interface{}, error) {
		val, err := o.getter.Get(key)
		return val, err
	})

	if err != nil {
		return nil, err
	}

	return val, nil
}

func (l *LRUCache) doGet(key string) (interface{}, error) {
	l.rwMutex.RLock()
	defer l.rwMutex.RUnlock()
	var (
		value *Node
		exist bool
	)
	if value, exist = l.valueMap[key]; !exist {
		return nil, ErrKeyNotExist
	}
	if value.val.isExpired() {
		return nil, ErrKeyExpired
	}
	l.list.moveToFront(value)
	return value.val.value, nil
}

func (l *LRUCache) Set(key string, val interface{}, opts ...OptionsFn) {
	o := getDefaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	node := NewNode(key, item{
		value:      val,
		expiration: time.Now().Add(o.expiration),
	})
	if int64(len(l.valueMap)) > l.maxSum && !l.isCleanIng() {
		l.Clean()
	}

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
