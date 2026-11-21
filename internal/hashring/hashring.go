package hashring

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"slices"
	"strconv"
	"sync"
)

type HashServer struct {
	ServerName string
	ServerIP   string
}

type HashRing struct {
	mx             sync.RWMutex
	ring           map[uint32]HashServer
	serverReplicas int
	sortedNodes    []uint32
}

func InitHashring(serverReplicas int) *HashRing {
	return &HashRing{
		ring:           make(map[uint32]HashServer),
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

		ring.ring[hash] = HashServer{
			ServerName: serverReplica,
			ServerIP:   serverIP,
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

func (ring *HashRing) FindKeyServer(key string) (HashServer, error) {
	ring.mx.RLock()
	defer ring.mx.RUnlock()

	if len(ring.sortedNodes) == 0 {
		return HashServer{}, errors.New("The ring is empty")
	}

	hash := generateHash(key)
	server, exists := ring.ring[hash]

	if !exists {
		index, _ := slices.BinarySearch(ring.sortedNodes, hash)

		if index == len(ring.sortedNodes) {
			index = 0
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
