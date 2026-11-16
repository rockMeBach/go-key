package store

import (
	"time"
)

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
