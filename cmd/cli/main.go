package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/google/shlex"
	"github.com/rockMeBach/gokey/internal/resp"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to go-key server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	fmt.Println("Welcome to go-key cli!")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("gokey>")

		if !scanner.Scan() {
			break
		}

		parts, err := shlex.Split(scanner.Text())
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid input:", err)
			continue
		}

		if len(parts) == 0 {
			continue
		}

		conn.Write(resp.EncodeRESP(parts))

		msg, err := resp.DecodeRESPClient(reader)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			fmt.Println(msg)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
