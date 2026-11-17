# Go-key: Cache store

Go-key is inspired by Redis and is written in Go. The C version is called C-key.

Key features and functions:

- Accepts GET, SET, DEL, EXPIRE commands
- Bi-directional communication using TCP
- Uses LRU cache-eviction policy
- Only string data is supported as of now
- Deletes expired keys using an active sweeper + passively whenever GET command is run
- Uses redis-like sampling technique to delete expired keys in the active sweeper
- Encodes and decodes commands and messages using RESP (Redis Serialization Protocol)

Features to be added:

- Data persistence
- Read from disk if the requested data is not found in the store
- Store write logs on disk

Command-line clients:

- Has its own client which can be run using cmd/cli/main.go
- You can also send commands via redis-cli
- Use valkey-cli if you're using it on Arch
