package main

import (
	"log"
	"os"
	"strings"

	"github.com/rockMeBach/gokey/internal/hashring"
	"github.com/rockMeBach/gokey/internal/proxy"
)

func main() {
	hr := hashring.InitHashring(1)

	STORES, set := os.LookupEnv("GOKEY_PROXY_STORES")

	if !set {
		log.Fatal("GOKEY_PROXY_STORES variable not set")
		return
	}

	stores := strings.Split(STORES, ",")

	if len(stores) == 0 {
		log.Fatal("No stores configured")
		return
	}

	for _, store := range stores {
		storeinfo := strings.Split(store, "-")
		name := storeinfo[0]
		address := storeinfo[1]

		err := hr.AddServer(name, address)

		if err != nil {
			panic(err)
		}
	}

	PORT, set := os.LookupEnv("GOKEY_PROXY_PORT")

	if !set {
		log.Fatal("GOKEY_PROXY_PORT variable not set")
		return
	}

	proxy.SetupProxyTCP(PORT, hr)
}
