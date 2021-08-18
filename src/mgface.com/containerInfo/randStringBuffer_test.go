package containerInfo

import (
	"fmt"
	"testing"
	"time"
)

func TestRandStringBuffer(t *testing.T) {
	for a := 0; a < 10; a++ {
		time.Sleep(1 * time.Second)
		fmt.Println(randStringBuffer(10))
	}
}
