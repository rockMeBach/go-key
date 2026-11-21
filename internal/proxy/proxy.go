package proxy

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/rockMeBach/gokey/internal/hashring"
	"github.com/rockMeBach/gokey/internal/resp"
)

func SetupProxyTCP(port string, ring *hashring.HashRing) {
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

		go handleConnection(conn, ring)
	}
}

func handleConnection(conn net.Conn, ring *hashring.HashRing) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		parts, err := resp.DecodeRESP(reader)
		if err != nil {
			return
		}

		if len(parts) < 2 {
			return
		}

		key := parts[1]
		keyServer, err := ring.FindKeyServer(key)

		fmt.Println("Key goes to", keyServer)

		if err != nil {
			return
		}

	}
}
