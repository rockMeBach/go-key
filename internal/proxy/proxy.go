package proxy

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/rockMeBach/gokey/internal/hashring"
	"github.com/rockMeBach/gokey/internal/resp"
)

type backend struct {
	conn   net.Conn
	reader *bufio.Reader
}

func SetupProxyTCP(address string, ring *hashring.HashRing) {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal("Error listening:", err)
	}

	defer listener.Close()

	fmt.Println("Now accepting connections on", address)

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
	backends := make(map[string]*backend)

	defer func() {
		for _, b := range backends {
			b.conn.Close()
		}
	}()

	for {
		parts, err := resp.DecodeRESP(reader)
		if err != nil {
			return
		}

		if len(parts) < 2 {
			conn.Write(resp.Error("wrong number of arguments"))
			continue
		}

		keyServer, err := ring.FindKeyServer(parts[1])
		if err != nil {
			conn.Write(resp.Error(err.Error()))
			continue
		}

		ip := keyServer.ServerIP
		b, exists := backends[ip]
		if !exists {
			c, err := net.Dial("tcp", ip)
			if err != nil {
				conn.Write(resp.Error("store unavailable"))
				continue
			}
			b = &backend{conn: c, reader: bufio.NewReader(c)}
			backends[ip] = b
		}

		reply, err := forward(b, parts)

		if err != nil {
			b.conn.Close()
			delete(backends, ip)
			conn.Write(resp.Error("store connection failed"))
			continue
		}

		if _, err := conn.Write(reply); err != nil {
			return
		}
	}
}

func forward(b *backend, parts []string) ([]byte, error) {
	if _, err := b.conn.Write(resp.EncodeRESP(parts)); err != nil {
		return nil, err
	}
	return resp.ReadRawReply(b.reader)
}
