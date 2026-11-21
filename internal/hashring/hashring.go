package hashring

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"slices"
	"strconv"
	"sync"
)

type Server struct {
	name string
	ip   string
}

type HashRing struct {
	mx             sync.RWMutex
	ring           map[uint32]Server
	serverReplicas int
	sortedNodes    []uint32
}

func InitHashring(serverReplicas int) *HashRing {
	return &HashRing{
		ring:           make(map[uint32]Server),
		serverReplicas: serverReplicas,
		sortedNodes:    make([]uint32, 0),
	}
}

func (ring *HashRing) AddServer(serverName string, serverIP string) error {
	ring.mx.Lock()
	defer ring.mx.Unlock()

	for i := range ring.serverReplicas {
		serverReplica := serverName + strconv.Itoa(i+1)
		hash := generateHash(serverReplica)

		_, exists := ring.ring[hash]

		if exists {
			return errors.New("Server with this name already exists")
		}

		ring.ring[hash] = Server{
			name: serverReplica,
			ip:   serverIP,
		}
		ring.sortedNodes = append(ring.sortedNodes, hash)
	}

	slices.Sort(ring.sortedNodes)
	return nil
}

func (ring *HashRing) RemoveServer(serverName string) error {
	ring.mx.Lock()
	defer ring.mx.Unlock()

	for i := range ring.serverReplicas {
		serverReplica := serverName + strconv.Itoa(i+1)
		hash := generateHash(serverReplica)

		_, exists := ring.ring[hash]

		if !exists {
			return errors.New("Server with this name doesn't exist")
		}

		delete(ring.ring, hash)

		index, _ := slices.BinarySearch(ring.sortedNodes, hash)
		ring.sortedNodes = slices.Delete(ring.sortedNodes, index, index+1)
	}

	return nil
}

func (ring *HashRing) FindKeyServer(key string) (Server, error) {
	ring.mx.RLock()
	defer ring.mx.RUnlock()

	if len(ring.sortedNodes) == 0 {
		return Server{}, errors.New("The ring is empty")
	}

	hash := generateHash(key)
	server, exists := ring.ring[hash]

	if !exists {
		index, _ := slices.BinarySearch(ring.sortedNodes, hash)

		if index == len(ring.sortedNodes) {
			index -= 1
		}

		closestServer := ring.ring[ring.sortedNodes[index]]
		return closestServer, nil
	}

	return server, nil
}

func generateHash(key string) uint32 {
	hash := md5.Sum([]byte(key))
	return binary.BigEndian.Uint32(hash[:])
}
