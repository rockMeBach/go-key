package store

import (
	"hash/fnv"
	"sync"
)

const SHARDS uint32 = 16
const MAX_CACHE_SIZE int = 30

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

func FindShard(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32() % SHARDS
}
