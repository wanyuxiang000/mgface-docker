package containerInfo

import (
	"fmt"
	"testing"
	"time"
)

func TestRandStrinByte(t *testing.T){
	for a:=0;a<10 ;a++  {
		time.Sleep(1*time.Second)
		fmt.Println(randStrinByte(10))
	}
}