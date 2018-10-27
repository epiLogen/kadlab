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

	file = fs.GetFile(key)

}

func (fs *FileSystem) GetRepublish(republishmin int) []File { //Returns the files that should be republished
	//fmt.Println("Getrepublish initierad")
	svar := []File{}
	now := time.Now()
	for i := 0; i < len(fs.files); i++ {
		fmt.Println("File diff time and comptime",now.Sub(fs.files[i].time), "Separering", time.Minute*time.Duration(republishmin-1))
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

func (fs *FileSystem) GetFile(key KademliaID) string { //Returns the content of the file
	fs.mtx.Lock()

	svar := ""
	pos := fs.GetPos(key)
	if pos != -1 {
		svar = fs.files[pos].data
	}
	fs.mtx.Unlock()

	fmt.Println("Getfile called key is", key.String())

	return svar
}

func (fs *FileSystem) GetPublisher(key KademliaID) string { //returns the publisher
	fs.mtx.Lock()
	svar := ""
	pos := fs.GetPos(key)
	if pos != -1 {
		svar = fs.files[pos].publisher
	}
	fs.mtx.Unlock()
	return svar
}

func (fs *FileSystem) Expired(key KademliaID) bool { //Kallas med låst mtx returnerar om en fil har expirat
	fmt.Println("Expired started")
	timen := time.Now()
	pos := fs.GetPos(key)
	if pos != -1 {
		fmt.Println("file exists")
		timen = fs.files[pos].time
	}
	now := time.Now()
	age := now.Sub(timen)
	fmt.Println("age is", age)
	if age > time.Minute*24*60 {
		return true
	} else {
		return false
	}
}

func (fs *FileSystem) Pinned(key KademliaID) bool {
	svar := false
	pos := fs.GetPos(key)
	if pos != -1 {
		svar = fs.files[pos].pin
	}
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

func (fs *FileSystem) Delete(key KademliaID) { //Kallas med låst mtx
	fmt.Println("Delete has been called")
	pos := fs.GetPos(key)
	if pos == -1 {
		fmt.Println("File does not exist length is", len(fs.files))
	} else{
		fmt.Println("Pos is not -1")
	}
	if !fs.Pinned(key) {
		fmt.Println("File exists and is not pinned")
		if (pos == 0) && (len(fs.files) == 1) {
			fmt.Println("case 1")
			fs.files = fs.files[:0]
			fmt.Println("deletion done")
		} else if pos == 0 {
			fs.files = fs.files[pos+1:]
		} else if pos == len(fs.files) {
			fs.files = fs.files[:pos-1]
		} else {
			fs.files = append(fs.files[:pos-1], fs.files[pos+1:]...)
		}
	} else{
		fmt.Println("File is pinned")
	}
}
