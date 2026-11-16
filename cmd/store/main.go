package main

import (
	"time"

	"github.com/rockMeBach/gokey/internal/server"
	"github.com/rockMeBach/gokey/internal/store"
)

func main() {
	shards := store.InitShards()

	stopSweeper := shards.StartSweeper(time.Second)
	defer stopSweeper()

	server.SetupTCP(":8080", shards)
}
