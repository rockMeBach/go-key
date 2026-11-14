package main

import (
	"fmt"
	"testing"
)

func TestEviction(t *testing.T) {
	shards := InitShards()

	// collect MAX_CACHE_SIZE+1 keys that all land in shard 0
	var keys []string
	for i := 0; len(keys) < MAX_CACHE_SIZE+1; i++ {
		k := fmt.Sprintf("key%d", i)
		if FindShard(k) == 0 {
			keys = append(keys, k)
		}
	}

	// insert them all: the last insert should evict keys[0]
	for _, k := range keys {
		shards.SetKey(k, "v")
	}

	if _, ok := shards.GetKey(keys[0]); ok {
		t.Errorf("oldest key %s should have been evicted", keys[0])
	}
	for _, k := range keys[1:] {
		if _, ok := shards.GetKey(k); !ok {
			t.Errorf("key %s should still exist", k)
		}
	}
}
