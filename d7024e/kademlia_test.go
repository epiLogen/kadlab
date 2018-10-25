package d7024e

import (
	"testing"
  "fmt"
  "strconv"
  "time"
)

func TestKademlia(t *testing.T) {

  n := 30

  contacts := make([]Contact, n)

  for i := 1; i <= n; i++ {
    intstring := strconv.Itoa(i)
    length := len(intstring)
    port := ""

    if length == 1 {
      port = "888" + intstring
    } else if length == 2 {
      port = "88" + intstring
    } else {
      port = "8" + intstring
    }

    contacts[i-1] = NewContact(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF" + port), "127.0.0.1:" + port)
  }

  nodes := make([]*Kademlia, n)

  for j := 0; j < n; j++ {
    nodes[j] = NewKademlia(contacts[j])
    go nodes[j].GetNetwork().Listen(contacts[j])
  }

  for k := 1; k < n; k++ {
    nodes[k].GetNetwork().GetRT().AddContact(contacts[0])
    go nodes[k].LookupContact(&contacts[k])
    time.Sleep(50*time.Millisecond)
  }

  nodes[0].Store("test")
	fmt.Println("STORE I KADTEST, SLEEPING...")
  time.Sleep(10*1000*time.Millisecond)
	fmt.Println("PRE LOOKUPDATA I KADTEST")
  //result := nodes[0].LookupData("test")
	fmt.Println("AFTER LOOKUPDATA I KADTEST")



	fmt.Println("hello from kadtest")
}
