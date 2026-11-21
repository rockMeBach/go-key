package main

import (
	"fmt"
	"strconv"

	"github.com/rockMeBach/gokey/internal/hashring"
)

func main() {
	hr := hashring.InitHashring(3)
	hr.AddServer("server1", "plaplapla")
	hr.AddServer("server2", "pliplipl")

	for i := range 5 {
		server, err := hr.FindKeyServer("key" + strconv.Itoa(i))

		if err == nil {
			fmt.Println("after adding", server)
		}
	}

	hr.RemoveServer("server2")

	for i := range 5 {
		server, err := hr.FindKeyServer("key" + strconv.Itoa(i))

		if err == nil {
			fmt.Println("after removing", server)
		}
	}
}
