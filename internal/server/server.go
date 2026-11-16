package server

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/rockMeBach/gokey/internal/resp"
	"github.com/rockMeBach/gokey/internal/store"
)

func SetupTCP(port string, shards *store.Shards) {
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatal("Error listening:", err)
	}

	defer listener.Close()

	fmt.Println("Now accepting connections on port", port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal("Error accepting connections:", err)
			continue
		}

		go handleConnection(conn, shards)
	}
}

func handleConnection(conn net.Conn, shards *store.Shards) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		parts, err := resp.DecodeRESP(reader)
		if err != nil {
			return
		}

		reply := dispatch(parts, shards)
		conn.Write(reply)
	}
}
