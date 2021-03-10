package local_cache

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	keys     []uint32
	hashMap  map[uint32]string
	multiple int
}

func New(multiple int, hash Hash) *Map {
	if hash == nil {
		hash = crc32.ChecksumIEEE
	}

	return &Map{
		hash:     hash,
		hashMap:  make(map[uint32]string),
		multiple: multiple,
	}
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.multiple; i++ {
			hashKey := m.hash([]byte(key + strconv.Itoa(i)))
			m.keys = append(m.keys, hashKey)
			m.hashMap[hashKey] = key
		}
	}

	sort.Slice(m.keys, func(i, j int) bool {
		return m.keys[i] < m.keys[j]
	})
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hashKey := m.hash([]byte(key))

	index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] > hashKey
	})

	return m.hashMap[m.keys[index%len(m.keys)]]
}
