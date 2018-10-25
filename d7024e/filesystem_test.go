package d7024e

import (
	"testing"
  "fmt"
  "time"
)

func TestFileSystem(t *testing.T) {


  contact := NewContact(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF" + "8001"), "127.0.0.1:" + "8001")
  node := NewKademlia(contact)
  key := NewKademliaID("test")

  pub := node.getFileSystem().GetPublisher(*key)
  expired := node.getFileSystem().Expired(*key)
  pinned := node.getFileSystem().Pinned(*key)
  repubfiles := node.getFileSystem().GetRepublish(2)
  node.getFileSystem().Pin(*key)
  node.getFileSystem().Unpin(*key)
  //node.getFileSystem().Delete(*key)




  time.Sleep(1*1000*time.Millisecond)


  fmt.Println("pub:", pub, "expired", expired, "pinned", pinned, "repubfiles", repubfiles)
	fmt.Println("hello from filesystest")
}
