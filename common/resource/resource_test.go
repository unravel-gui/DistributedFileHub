package resource

import (
	"fmt"
	"testing"
)

func TestNewResourceStatus(t *testing.T) {
	x := NewResourceStatus()
	fmt.Println(x)
}
