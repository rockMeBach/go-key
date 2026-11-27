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
	conn, err := net.Dial("tcp", ":6379")
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

		if len(parts) == 1 && parts[0] == "exit" {
			fmt.Println("quitting")
			return
		}

		conn.Write(resp.EncodeRESP(parts))

		msg, err := resp.DecodeRESPClient(reader)

		// show the output on the cmd line
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
