package d7024e

import (
	"testing"
  "fmt"
)

func TestKademliaid(t *testing.T) {

	k1 := NewKademliaIDnp("hej")
  k2 := NewRandomKademliaID()
	b1 := k1.Less(&k2)
	b2 := k1.Equals(&k2)
  fmt.Println(k1, b1)
  fmt.Println(k2, b2)

	//

}
