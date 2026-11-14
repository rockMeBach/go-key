package main

import (
	"fmt"
	"hash/fnv"
	"sync"
)

const SHARDS uint32 = 16

type SafeStore struct {
	mx    sync.Mutex
	store map[string]string
}

type Shards struct {
	storeArr []SafeStore
}

func InitShards() *Shards {
	shard := &Shards{
		storeArr: make([]SafeStore, SHARDS),
	}

	for i := range shard.storeArr {
		shard.storeArr[i].store = make(map[string]string)
	}

	return shard
}

func (shards *Shards) SetKey(key string, value string) {
	sh := &shards.storeArr[FindShard(key)]

	sh.mx.Lock()
	defer sh.mx.Unlock()

	sh.store[key] = value
}

func (shards *Shards) GetKey(key string) (string, bool) {
	sh := &shards.storeArr[FindShard(key)]

	sh.mx.Lock()
	defer sh.mx.Unlock()

	value, ok := sh.store[key]
	return value, ok
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
