package store

type Node struct {
	key, value string
	expiresAt  int64
	prev, next *Node
}

func (store *SafeStore) unlink(node *Node) {
	// Just unlink the node, nothing else
	// Update the DLL after unlinking

	if node.prev != nil {
		node.prev.next = node.next
	} else {
		store.lru = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		store.mru = node.prev
	}

	node.next, node.prev = nil, nil
}

func (store *SafeStore) updateMRU(node *Node) {
	node.next = nil
	node.prev = store.mru

	if store.mru != nil {
		store.mru.next = node
	} else {
		store.lru = node
	}

	store.mru = node
}
