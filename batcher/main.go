package main

import (
	"fmt"
	"time"
)

func main() {
	//Count := 12345
	//batcher := NewBatcher(16, 100, time.Second)
	//if batcher == nil {
	//	panic("failed to create batcher")
	//}
	//
	//batcher.Start()
	//
	//go func() {
	//	for i := 0; i < Count; i++ {
	//		batcher.In() <- i
	//		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	//	}
	//	batcher.Stop()
	//}()
	//
	//for {
	//	select {
	//	case <-batcher.stop:
	//		fmt.Println("stop")
	//		return
	//	case datas := <-batcher.Out():
	//		fmt.Printf("resv size:%d\n", len(datas))
	//	}
	//}

	timer := time.NewTimer(time.Millisecond)
	time.Sleep(time.Second)
	fmt.Printf("stop:%v\n", timer.Stop())
	fmt.Printf("reset:%v\n", timer.Reset(time.Second))

	select {
	case <-timer.C:
		fmt.Println("读到了1")
	default:
		fmt.Println("没读到2")
	}
	time.Sleep(5 * time.Second)
	select {
	case <-timer.C:
		fmt.Println("读到了2")
	default:
		fmt.Println("没读到2")
	}

}
