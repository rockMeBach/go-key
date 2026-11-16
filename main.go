package main

import (
	"hash/fnv"
	"sync"
	"time"
)

const SHARDS uint32 = 16
const MAX_CACHE_SIZE int = 30

type Node struct {
	key, value string
	expiresAt  int64
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
			sh.size -= 1
			old := sh.lru
			sh.unlink(old)
			delete(sh.store, old.key)
		}

	} else {
		node.value = value
		node.expiresAt = 0
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
		if node.expiresAt != 0 && time.Now().UnixMilli() > node.expiresAt {
			sh.unlink(node)
			delete(sh.store, key)
			sh.size -= 1

			return "", false
		}

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

func (shards *Shards) ExpireKey(key string, seconds uint) bool {
	sh := &shards.storeArr[FindShard(key)]

	sh.mx.Lock()
	defer sh.mx.Unlock()

	node, ok := sh.store[key]
	now := time.Now().UnixMilli()

	if ok {
		if node.expiresAt != 0 && now > node.expiresAt {
			return false
		}

		node.expiresAt = now + int64(seconds)*1000
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

func (safeStore *SafeStore) sweepStore(passCap int, batchSize int) {
	// Redis-like sampling
	for range passCap {
		safeStore.mx.Lock()

		count := 0
		totalWithTTL := 0
		totalExpired := 0
		now := time.Now().UnixMilli()

		for key, node := range safeStore.store {
			if node.expiresAt != 0 {
				totalWithTTL += 1

				if now > node.expiresAt {
					totalExpired += 1
					safeStore.unlink(node)
					safeStore.size -= 1
					delete(safeStore.store, key)
				}
			}

			count += 1

			if count == batchSize {
				break
			}
		}

		safeStore.mx.Unlock()

		if totalWithTTL == 0 {
			return
		}

		if float64(totalExpired)/float64(totalWithTTL) <= 0.25 {
			return
		}
	}
}

func (shards *Shards) startSweeper(interval time.Duration) (stop func()) {
	done := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for i := range shards.storeArr {
					shards.storeArr[i].sweepStore(10, 20)
				}
			}
		}
	}()

	return func() { close(done) }
}

func main() {

}
