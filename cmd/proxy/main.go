package main

import (
	"github.com/rockMeBach/gokey/internal/hashring"
	"github.com/rockMeBach/gokey/internal/proxy"
)

func main() {
	hr := hashring.InitHashring(1)

	hr.AddServer("server1", ":8080")
	hr.AddServer("server2", ":8080")

	proxy.SetupProxyTCP(":8000", hr)
}
