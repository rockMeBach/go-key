package store

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

func fillExpired(sh *SafeStore, n int) {
	sh.store = make(map[string]*Node)
	past := time.Now().Add(-time.Minute).UnixMilli()

	for i := range n {
		key := fmt.Sprintf("key%d", i)
		node := &Node{key: key, value: "v", expiresAt: past}
		sh.store[key] = node
		sh.updateMRU(node)
		sh.size++
	}
}

func TestSweeperRepeatPasses(t *testing.T) {
	sh := &SafeStore{}
	fillExpired(sh, 100)

	// every sampled key is expired, so the ratio stays at 100%
	// and it keeps repeating: 5 passes of 20 clear all 100
	sh.sweepStore(10, 20)

	if len(sh.store) != 0 {
		t.Errorf("expected empty map, got %d keys", len(sh.store))
	}
	if sh.size != 0 {
		t.Errorf("expected size 0, got %d", sh.size)
	}
	if sh.lru != nil || sh.mru != nil {
		t.Error("list should be empty after all nodes are unlinked")
	}
}

func TestSweeperPassCap(t *testing.T) {
	sh := &SafeStore{}
	fillExpired(sh, 100)

	// only 2 passes allowed, so exactly 2 x 20 = 40 removed
	sh.sweepStore(2, 20)

	if len(sh.store) != 60 {
		t.Errorf("expected 60 keys left, got %d", len(sh.store))
	}
	if sh.size != 60 {
		t.Errorf("expected size 60, got %d", sh.size)
	}
}

func TestSweeperMixed(t *testing.T) {
	sh := &SafeStore{store: make(map[string]*Node)}
	now := time.Now()
	past := now.Add(-time.Minute).UnixMilli()
	future := now.Add(time.Hour).UnixMilli()

	add := func(key string, expiresAt int64) {
		node := &Node{key: key, value: "v", expiresAt: expiresAt}
		sh.store[key] = node
		sh.updateMRU(node)
		sh.size++
	}

	for i := range 10 {
		add(fmt.Sprintf("expired%d", i), past)
		add(fmt.Sprintf("live%d", i), future)
		add(fmt.Sprintf("nottl%d", i), 0)
	}

	// batch covers the whole map, so one pass sees every key
	sh.sweepStore(1, 100)

	if len(sh.store) != 20 {
		t.Errorf("expected 20 keys left, got %d", len(sh.store))
	}
	if sh.size != 20 {
		t.Errorf("expected size 20, got %d", sh.size)
	}

	for i := range 10 {
		if _, ok := sh.store[fmt.Sprintf("expired%d", i)]; ok {
			t.Errorf("expired%d should be gone", i)
		}
		if _, ok := sh.store[fmt.Sprintf("live%d", i)]; !ok {
			t.Errorf("live%d should remain", i)
		}
		if _, ok := sh.store[fmt.Sprintf("nottl%d", i)]; !ok {
			t.Errorf("nottl%d should remain", i)
		}
	}

	// the list must still be consistent: walk it and count
	n := 0
	for node := sh.lru; node != nil; node = node.next {
		n++
	}
	if n != 20 {
		t.Errorf("list has %d nodes, expected 20", n)
	}
}
