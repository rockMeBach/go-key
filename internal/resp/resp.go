package resp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
)

const MAX_BULK int = 1000

func SimpleString(s string) []byte { return []byte("+" + s + "\r\n") }
func Error(msg string) []byte      { return []byte("-ERR " + msg + "\r\n") }
func Integer(n int) []byte         { return []byte(":" + strconv.Itoa(n) + "\r\n") }
func Null() []byte                 { return []byte("$-1\r\n") }
func BulkString(s string) []byte {
	return []byte("$" + strconv.Itoa(len(s)) + "\r\n" + s + "\r\n")
}

func EncodeRESP(command []string) []byte {
	var buffer bytes.Buffer

	buffer.WriteString("*")
	buffer.WriteString(strconv.Itoa(len(command)))
	buffer.WriteString("\r\n")

	for _, part := range command {
		buffer.WriteString("$")
		buffer.WriteString(strconv.Itoa(len(part)))
		buffer.WriteString("\r\n")
		buffer.WriteString(part)
		buffer.WriteString("\r\n")
	}

	return buffer.Bytes()
}

func DecodeRESP(reader *bufio.Reader) ([]string, error) {
	peek, err := reader.Peek(1)

	if err != nil {
		return nil, err
	}

	if peek[0] != '*' {
		return nil, errors.New("protocol error, the first character should be a *")
	}

	line, err := reader.ReadString('\n')

	if err != nil {
		return nil, err
	}

	N_string := strings.TrimSuffix(line[1:], "\r\n")

	N, err := strconv.Atoi(N_string)

	if err != nil || N < 0 || N > MAX_BULK {
		return nil, errors.New("protocol error, invalid bulk length")
	}

	parts := make([]string, 0, N)

	for range N {
		line, err = reader.ReadString('\n')

		if line[0] != '$' {
			return nil, errors.New("protocol error, $ not found")
		}

		L, err := strconv.Atoi(strings.TrimSuffix(line[1:], "\r\n"))

		if err != nil || L < 0 || L > MAX_BULK {
			return nil, errors.New("protocol error, invalid bulk length")
		}

		// Read exactly L bytes along with the \r and \n (L+2)
		buf := make([]byte, L+2)

		if _, err := io.ReadFull(reader, buf); err != nil {
			return nil, err
		}

		if buf[L] != '\r' || buf[L+1] != '\n' {
			return nil, errors.New("protocol error, missing CRLF after bulk string")
		}

		parts = append(parts, string(buf[:L]))
	}

	return parts, nil
}
