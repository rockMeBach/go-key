package main

import (
	"fmt"

	store "github.com/rockMeBach/gokey/internal/store"
)

func main() {
	shards := store.InitShards()
	shards.SetKey("testkey", "v")
	fmt.Println(shards.GetKey("testkey"))
}
