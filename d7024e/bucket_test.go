package d7024e

import (
	"testing"
  "fmt"
)

func TestBucket(t *testing.T) {

	//me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	//rt := NewRoutingTable(me)

	c1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	//c2 := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002")
	//c3 := NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002")
	//c4 := NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002")
	//c5 := NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002")
	//c6 := NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002")

	//rt.AddContact(c1)
	//rt.AddContact(c2)
	//rt.AddContact(c3)
	//rt.AddContact(c4)
	//rt.AddContact(c5)
	//rt.AddContact(c6)

  bucket := newBucket()
  bucket.AddContact(c1)
  c1inbucket := bucket.ContactinBucket(c1)
  bucket.RemoveContact(c1)
  c1inbucket = bucket.ContactinBucket(c1)
  size := bucket.Len()

  fmt.Println("c1inbucket:", c1inbucket, "size", size)

}
