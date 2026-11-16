package main

import (
	"fmt"
	"testing"
	"time"
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

func TestDeletion(t *testing.T) {
	shards := InitShards()

	key := "key1"

	shards.SetKey(key, "v")

	if _, ok := shards.GetKey(key); !ok {
		t.Error("key was expected to be set")
	}

	shards.DelKey(key)

	if _, ok := shards.GetKey(key); ok {
		t.Error("key was expected to be deleted")
	}

	if ok := shards.DelKey(key); ok {
		t.Error("DelKey should have returned false")
	}
}

func TestExpiration(t *testing.T) {
	// could use smaller test times
	shards := InitShards()

	key := "key1"

	t.Log("Testing to check whether GetKey detects expired keys...")

	shards.SetKey(key, "v")

	if !shards.ExpireKey(key, 5) {
		t.Fatal("ExpireKey should have found the key")
	}

	if _, ok := shards.GetKey(key); !ok {
		t.Error("GetKey should have returned true")
	}

	time.Sleep(time.Duration(6) * time.Second)

	if _, ok := shards.GetKey(key); ok {
		t.Error("GetKey should have returned false")
	}

	t.Log("Testing to check whether ExpireKey detects expired keys...")

	shards.SetKey(key, "v")
	shards.ExpireKey(key, 3)

	time.Sleep(time.Duration(4) * time.Second)

	if ok := shards.ExpireKey(key, 5); ok {
		t.Error("ExpireKey should have returned false")
	}
}

func TestSweeper(t *testing.T) {
	shards := InitShards()
	stop := shards.startSweeper(100 * time.Millisecond)
	defer stop()

	key := "key1"
	shards.SetKey(key, "v")
	if !shards.ExpireKey(key, 1) {
		t.Fatal("ExpireKey should have found the key")
	}

	time.Sleep(1500 * time.Millisecond)

	// inspect the map directly, NOT through GetKey
	sh := &shards.storeArr[FindShard(key)]
	sh.mx.Lock()
	defer sh.mx.Unlock()

	if _, exists := sh.store[key]; exists {
		t.Error("sweeper should have removed the expired key from the map")
	}
	if sh.size != 0 {
		t.Errorf("size should be 0, got %d", sh.size)
	}
}
