package main

import (
	"log"
	"os"
	"time"

	"github.com/rockMeBach/gokey/internal/server"
	"github.com/rockMeBach/gokey/internal/store"
)

func main() {
	shards := store.InitShards()

	stopSweeper := shards.StartSweeper(time.Second)
	defer stopSweeper()

	PORT, err := os.LookupEnv("GOKEY_STORE_PORT")

	if !err {
		log.Fatal("GOKEY_STORE_PORT variable not set")
		return
	}

	server.SetupTCP(PORT, shards)
}
