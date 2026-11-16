package server

import (
	"errors"
	"strconv"
	"strings"

	"github.com/rockMeBach/gokey/internal/store"
)

func dispatch(parts []string, shards *store.Shards) (string, string, error) {
	if len(parts) == 0 {
		return "", "", errors.New("empty parts")
	}

	command := strings.ToUpper(parts[0])
	miss_arg_err := errors.New("missing arguments")

	switch command {
	case "SET":
		if len(parts) != 3 {
			return "", "", miss_arg_err
		}

		key := parts[1]
		value := parts[2]
		shards.SetKey(key, value)
		return "key set", "", nil
	case "GET":
		if len(parts) != 2 {
			return "", "", miss_arg_err
		}

		key := parts[1]

		val, ok := shards.GetKey(key)

		if !ok {
			return "", "", errors.New("key not found")
		}

		return "key retrieved", val, nil
	case "DEL":
		if len(parts) != 2 {
			return "", "", miss_arg_err
		}

		key := parts[1]

		if shards.DelKey(key) {
			return "key deleted", "", nil
		}

		return "", "", errors.New("key not found")
	case "EXPIRE":
		if len(parts) != 3 {
			return "", "", miss_arg_err
		}

		key := parts[1]
		seconds := parts[2]

		seconds_int, err := strconv.Atoi(seconds)

		if err != nil {
			return "", "", errors.New("error setting expiry")
		}

		if shards.ExpireKey(key, seconds_int) {
			return "expiry set", "", nil
		}

		return "", "", errors.New("key has expired")
	}

	return "", "", errors.New("unknown command")
}
