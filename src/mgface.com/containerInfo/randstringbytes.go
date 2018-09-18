package containerInfo

import (
	"math/rand"
	"time"
)

func randStrinByte(n int) string{
	letterBytes:="0123456789"
	b:=make([]byte,n)
	rand.Seed(time.Now().UnixNano())
	for i:=range b{
		b[i]=letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

