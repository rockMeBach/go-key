# Go-key: Cache store

Go-key is inspired by Redis and is written in Go and C. The C version is called C-key.

Key features:

- Accepts GET, SET, DEL, EXPIRE commands
- Bi-directional communication using TCP
- Uses LRU cache-eviction policy
- Only string data is supported as of now

Features to be added:

- Data persistence
- Read from disk if the requested data is not found in the store
- Store write logs on disk

Clients:

- Has its own client
- You can also send commands via redis-cli
- Use valkey-cli if you're using it on Arch
