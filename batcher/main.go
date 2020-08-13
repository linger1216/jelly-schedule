package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	Count := 12345
	batcher := NewBatcher(16, 100, time.Second)
	if batcher == nil {
		panic("failed to create batcher")
	}

	batcher.Start()

	go func() {
		for i := 0; i < Count; i++ {
			batcher.In() <- i
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
		}
		batcher.Stop()
	}()

	for {
		select {
		case <-batcher.stop:
			fmt.Println("stop")
			return
		case datas := <-batcher.Out():
			fmt.Printf("resv size:%d\n", len(datas))
		}
	}
}
