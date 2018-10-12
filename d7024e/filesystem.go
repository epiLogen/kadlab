package d7024e

import (
	"fmt"
	"sync"
	"time"
)

type FileSystem struct {
	publishers map[KademliaID]string
	data       map[KademliaID]string
	times      map[KademliaID]time.Time
	pin        map[KademliaID]bool
	mtx        *sync.Mutex
}

func NewFileSystem() FileSystem {
	var filesystem FileSystem
	filesystem.publishers = make(map[KademliaID]string)
	filesystem.data = make(map[KademliaID]string)
	filesystem.times = make(map[KademliaID]time.Time)
	filesystem.pin = make(map[KademliaID]bool)
	filesystem.mtx = &sync.Mutex{}
	return filesystem
}
