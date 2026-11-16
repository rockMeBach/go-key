package server

import (
	"strconv"
	"strings"

	"github.com/rockMeBach/gokey/internal/resp"
	"github.com/rockMeBach/gokey/internal/store"
)

func dispatch(parts []string, shards *store.Shards) []byte {
	if len(parts) == 0 {
		return resp.Error("empty command")
	}

	command := strings.ToUpper(parts[0])
	argErr := resp.Error("wrong number of arguments for '" + strings.ToLower(command) + "' command")

	switch command {
	case "SET":
		if len(parts) != 3 {
			return argErr
		}
		shards.SetKey(parts[1], parts[2])
		return resp.SimpleString("OK")

	case "GET":
		if len(parts) != 2 {
			return argErr
		}
		val, ok := shards.GetKey(parts[1])
		if !ok {
			return resp.Null()
		}
		return resp.BulkString(val)

	case "DEL":
		if len(parts) != 2 {
			return argErr
		}
		if shards.DelKey(parts[1]) {
			return resp.Integer(1)
		}
		return resp.Integer(0)

	case "EXPIRE":
		if len(parts) != 3 {
			return argErr
		}
		seconds, err := strconv.Atoi(parts[2])
		if err != nil || seconds < 0 {
			return resp.Error("value is not an integer or out of range")
		}
		if shards.ExpireKey(parts[1], seconds) {
			return resp.Integer(1)
		}
		return resp.Integer(0)
	}

	return resp.Error("unknown command '" + parts[0] + "'")
}
