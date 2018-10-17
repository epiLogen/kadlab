package d7024e

import (
	"sync"
	"time"
	"fmt"
)

type FileSystem struct {
//	publishers map[*KademliaID]string
//	data       map[*KademliaID]string
//	times      map[*KademliaID]time.Time
//	pin        map[*KademliaID]bool
	files 		[]File
	mtx        *sync.Mutex
}

type File struct {
	key KademliaID
	publisher string
	data string
	time time.Time
	pin bool
}

func NewFileSystem() FileSystem {
	var fs FileSystem
	fs.files = []File{}
//	fs.publishers = make(map[*KademliaID]string)
//	fs.data = make(map[*KademliaID]string)
//	fs.times = make(map[*KademliaID]time.Time)
//	fs.pin = make(map[*KademliaID]bool)
	fs.mtx = &sync.Mutex{}
	return fs
}

// func NewFile() File {
// 	var f File
// 	return fs
// }

func (fs *FileSystem) Store(key KademliaID, file string, publisher string) {
	fs.mtx.Lock()

	// fs.publishers[key] = publisher
	// fs.data[key] = file
	// fs.times[key] = time.Now()
	// fs.pin[key] = false

	// pos := len (fs.files) -1
	// if len (fs.files) == 0 {
	// 	pos = 0
	// }
	filen := File{}
	filen.publisher = publisher
	filen.data = file
	filen.key = key
	filen.time = time.Now()
	filen.pin = false
	fs.files = append(fs.files, filen)

	fs.mtx.Unlock()

	fmt.Println("File stored")
	file = fs.GetFile(key)
	fmt.Println("File is", file)

}

//returns the position of the searched file. Called with locked mutex
func (fs *FileSystem) GetPos(key KademliaID) int {
	svar := -1
	for i, file := range fs.files {
	if file.key == key {
		svar = i
		}
	}
	return svar
}

func (fs *FileSystem) GetFile(key KademliaID) string {
	fs.mtx.Lock()
//	size := len(fs.data)
//	svar := fs.data[key]

	svar := ""
	pos := fs.GetPos(key)
	if pos != -1 {
		svar = fs.files[pos].data
	}
	fs.mtx.Unlock()

	fmt.Println("Getfile called size is, key is, svar is", len(fs.files), key.String(), svar)

	//Print map
	// fmt.Println("GPrinting map")
	// for k, v := range fs.data {
  //   fmt.Printf("key[%s] value[%s]\n", k, v)
	// }
	return svar
}

func (fs *FileSystem) GetPublisher(key KademliaID) string {
	fs.mtx.Lock()
	svar := ""
	pos := fs.GetPos(key)
	if pos != -1 {
		svar = fs.files[pos].publisher
	}
	fs.mtx.Unlock()
	return svar
}

func (fs *FileSystem) Expired(key KademliaID) bool {
	fs.mtx.Lock()

	timen := time.Now()
	pos := fs.GetPos(key)
	if pos != -1 {
		timen = fs.files[pos].time
	}

	fs.mtx.Unlock()

	now := time.Now()
	age := now.Sub(timen)
	return age > time.Hour*24
}

func (fs *FileSystem) Pinned(key KademliaID) bool {
	fs.mtx.Lock()

	svar := false
	pos := fs.GetPos(key)
	if pos != -1 {
		svar = fs.files[pos].pin
	}

	fs.mtx.Unlock()
	return svar
}

func (fs *FileSystem) Pin(key KademliaID) {
	fs.mtx.Lock()
	pos := fs.GetPos(key)
	if pos != -1 {
		fs.files[pos].pin = true
	}
	fs.mtx.Unlock()
}

func (fs *FileSystem) Unpin(key KademliaID) {
	fs.mtx.Lock()
	pos := fs.GetPos(key)
	if pos != -1 {
		fs.files[pos].pin = false
	}
	fs.mtx.Unlock()
}

func (fs *FileSystem) Delete(key KademliaID) {
	fs.mtx.Lock()
	pos := fs.GetPos(key)
	if !fs.Pinned(key) {
		if (pos == 0) && (len(fs.files) == 1) {
			fs.files = []File{}
		} else if pos == 0 {
			fs.files = fs.files[pos+1:]
		} else if pos == len(fs.files) {
			fs.files = fs.files[:pos-1]
		} else {
			fs.files = append(fs.files[:pos-1], fs.files[pos+1:]...)
		}
	}
	fs.mtx.Unlock()
}
