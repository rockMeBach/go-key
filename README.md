# Go-key: Cache store

Go-key is inspired by Redis and is written in Go and C. The C version is called C-key.

Key features:

- Communication via TCP
- Accepts GET, SET commands
- Uses LRU cache-eviction policy
- Only string data is supported in version 1.0
- Reads from disk if the requested data is not found in the memory
- Write logs stored on disk