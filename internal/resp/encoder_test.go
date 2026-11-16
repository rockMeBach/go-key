package resp

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestEncoderDecoder(t *testing.T) {
	command := "SET hello hi"
	respBytes := EncodeRESP(strings.Split(command, " "))
	t.Logf("%q\n", respBytes)

	reader := bufio.NewReader(bytes.NewReader(respBytes))
	parts, err := DecodeRESP(reader)

	if err != nil {
		t.Error("error decoding")
		return
	}

	t.Log(parts)
}
