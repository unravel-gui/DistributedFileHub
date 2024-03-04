package temp

import (
	"fmt"
	"os"
	"testing"
)

func Test_RenamePath(t *testing.T) {
	f, err := os.Create("test.txt")
	if err != nil {
		return
	}
	err = os.Remove("test.txt")
	if err != nil {
		fmt.Println("Error moving file:", err)
		return
	}
	f.Close()
	err = os.Rename("test.txt", "ttt.txt")
	if err != nil {
		fmt.Println("Error moving file:", err)
		return
	}

}
