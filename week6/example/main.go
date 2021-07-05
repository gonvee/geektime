package main

import (
	"fmt"
	"geektime/week6"
	"math/rand"
	"time"
)

func main() {
	fixedSeconds := 5 // 5秒宽度窗口
	c := week6.NewCounter(fixedSeconds)

	go func() {
		for {
			c.Inc()
			time.Sleep(time.Duration(rand.Intn(30)) * time.Millisecond)
		}
	}()

	for {
		time.Sleep(time.Second)
		fmt.Printf("last %d seconds: %d\n", fixedSeconds, c.Count())
	}

}
