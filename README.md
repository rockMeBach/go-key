# Go-key: Cache store

Go-key is inspired by Redis and is written in Go and C. The C version is called C-key.

Key features:

- Accepts GET, SET, DEL, EXPIRE TCP commands
- Uses an LRU cache-eviction policy
- Supports lists and sets, not just strings
- Reads from disk if the requested data is not found in the memory
