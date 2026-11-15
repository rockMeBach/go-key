package main

import (
	"fmt"
	"hash/fnv"
	"sync"
)

const SHARDS uint32 = 16
const MAX_CACHE_SIZE int = 30

type Node struct {
	key, value string
	prev, next *Node
}

type SafeStore struct {
	mx    sync.Mutex
	store map[string]*Node
	size  int
	mru   *Node
	lru   *Node
}

type Shards struct {
	storeArr []SafeStore
}

func InitShards() *Shards {
	shard := &Shards{
		storeArr: make([]SafeStore, SHARDS),
	}

	for i := range shard.storeArr {
		shard.storeArr[i].store = make(map[string]*Node)
	}

	return shard
}

func (store *SafeStore) unlink(node *Node) {
	// Just unlink the node, nothing else
	// Update the DLL after unlinking

	if node.prev != nil {
		node.prev.next = node.next
	} else {
		store.lru = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		store.mru = node.prev
	}

	node.next, node.prev = nil, nil
}

func (store *SafeStore) updateMRU(node *Node) {
	node.next = nil
	node.prev = store.mru

	if store.mru != nil {
		store.mru.next = node
	} else {
		store.lru = node
	}

	store.mru = node
}

func (shards *Shards) SetKey(key string, value string) {
	sh := &shards.storeArr[FindShard(key)]

	sh.mx.Lock()
	defer sh.mx.Unlock()

	node, exists := sh.store[key]

	if !exists {
		node = &Node{
			key:   key,
			value: value,
		}

		sh.store[key] = node
		sh.size += 1
		sh.updateMRU(node)

		if sh.size > MAX_CACHE_SIZE {
			sh.size = MAX_CACHE_SIZE
			old := sh.lru
			sh.unlink(old)
			delete(sh.store, old.key)
		}

	} else {
		node.value = value
		sh.unlink(node)
		sh.updateMRU(node)
	}
}

func (shards *Shards) GetKey(key string) (string, bool) {
	sh := &shards.storeArr[FindShard(key)]

	sh.mx.Lock()
	defer sh.mx.Unlock()

	node, ok := sh.store[key]

	if ok {
		sh.unlink(node)
		sh.updateMRU(node)

		return node.value, ok
	}

	return "", ok
}

func (shards *Shards) DelKey(key string) bool {
	sh := &shards.storeArr[FindShard(key)]

	sh.mx.Lock()
	defer sh.mx.Unlock()

	node, ok := sh.store[key]

	if ok {
		sh.unlink(node)
		delete(sh.store, key)
		sh.size -= 1
	} else {
		return false
	}

	return true
}

func FindShard(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32() % SHARDS
}

func main() {
	shards := InitShards()
	shards.SetKey("theme", "light")
	fmt.Println(shards.GetKey("theme"))
}
