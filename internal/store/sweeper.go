package store

import "time"

func (safeStore *SafeStore) sweepStore(passCap int, batchSize int) {
	// Redis-like sampling
	for range passCap {
		safeStore.mx.Lock()

		count := 0
		totalWithTTL := 0
		totalExpired := 0
		now := time.Now().UnixMilli()

		for key, node := range safeStore.store {
			if node.expiresAt != 0 {
				totalWithTTL += 1

				if now > node.expiresAt {
					totalExpired += 1
					safeStore.unlink(node)
					safeStore.size -= 1
					delete(safeStore.store, key)
				}
			}

			count += 1

			if count == batchSize {
				break
			}
		}

		safeStore.mx.Unlock()

		if totalWithTTL == 0 {
			return
		}

		if float64(totalExpired)/float64(totalWithTTL) <= 0.25 {
			return
		}
	}
}

func (shards *Shards) StartSweeper(interval time.Duration) (stop func()) {
	done := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for i := range shards.storeArr {
					shards.storeArr[i].sweepStore(10, 20)
				}
			}
		}
	}()

	return func() { close(done) }
}
