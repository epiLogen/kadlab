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
	var fs FileSystem
	fs.publishers = make(map[KademliaID]string)
	fs.data = make(map[KademliaID]string)
	fs.times = make(map[KademliaID]time.Time)
	fs.pin = make(map[KademliaID]bool)
	fs.mtx = &sync.Mutex{}
	return fs
}

func (fs *FileSystem) GetFile(key *KademliaID) string {
	fs.mtx.Lock()
	defer fs.mtx.Unlock()
	return fs.data[*key]
}

func (fs *FileSystem) GetPublisher(key *KademliaID) string {
	fs.mtx.Lock()
	defer fs.mtx.Unlock()
	return fs.publishers[*key]
}

func (fs *FileSystem) Expired(key *KademliaID) bool {
	fs.mtx.Lock()
	defer fs.mtx.Unlock()
	now := time.Now()
	age := now.Sub(fs.times[*key])
	return age > time.Hour*24
}

func (fs *FileSystem) Pinned(key *KademliaID) bool {
	fs.mtx.Lock()
	defer fs.mtx.Unlock()
	return fs.pin[*key]
}

func (fs *FileSystem) Pin(key *KademliaID) {
	fs.mtx.Lock()
	defer fs.mtx.Unlock()
	fs.pin[*key] = true
}

func (fs *FileSystem) Unpin(key *KademliaID) {
	fs.mtx.Lock()
	defer fs.mtx.Unlock()
	fs.pin[*key] = false
}

func (fs *FileSystem) Delete(key *KademliaID) {
	fs.mtx.Lock()
	defer fs.mtx.Unlock()
	if !fs.Pinned(key) {
		delete(fs.data, *key)
		delete(fs.publishers, *key)
		delete(fs.times, *key)
		delete(fs.pin, *key)
	}
}
