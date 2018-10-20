package d7024e

import (
	"sync"
	"time"
	"fmt"
)

type FileSystem struct {
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
	fs.mtx = &sync.Mutex{}
	return fs
}

func (fs *FileSystem) Store(key KademliaID, file string, publisher string) {
	fs.mtx.Lock()

	pos := fs.GetPos(key)
	if pos != -1 {
		fs.files[pos].time = time.Now()
		fmt.Println("Republish registered")
	} else {
		filen := File{}
		filen.publisher = publisher
		filen.data = file
		filen.key = key
		filen.time = time.Now()
		filen.pin = false
		fs.files = append(fs.files, filen)
		fmt.Println("Store registered")
	}

	fs.mtx.Unlock()

	//file = fs.GetFile(key)

}

func (fs *FileSystem) GetRepublish(republishmin int) []File {
	fmt.Println("Getrepublish initierad")
	svar := []File{}
	now := time.Now()
	for i := 0; i < len(fs.files); i++ {
		fmt.Println("File diff time and comptime",now.Sub(fs.files[i].time), "Separering", time.Minute*time.Duration(republishmin-2))
		if now.Sub(fs.files[i].time) > time.Minute*time.Duration(republishmin-1) {
			svar = append(svar, fs.files[i])
		}
	}
	fmt.Println("Längden på svaret är", len(svar))
	return svar
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

	svar := ""
	pos := fs.GetPos(key)
	if pos != -1 {
		svar = fs.files[pos].data
	}
	fs.mtx.Unlock()

	fmt.Println("Getfile called size is, key is, svar is", len(fs.files), key.String(), svar)

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
