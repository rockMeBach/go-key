package main

import (
	"fmt"
	"strconv"

	"github.com/rockMeBach/gokey/internal/hashring"
)

func main() {
	hr := hashring.InitHashring(3)
	hr.AddServer("server1")
	hr.AddServer("server2")

	for i := range 20 {
		server, err := hr.FindKeyServer("key" + strconv.Itoa(i))

		if err == nil {
			fmt.Println(server)
		}
	}
}
