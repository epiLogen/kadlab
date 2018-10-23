package d7024e

import (
	"testing"
  "fmt"
)

func TestKademliaid(t *testing.T) {

	k1 := NewKademliaIDnp("hej")
  k2 := NewRandomKademliaID()
  fmt.Println(k1)
  fmt.Println(k2)

}
