package week6

import (
	"container/ring"
	"sync"
	"time"
)

type Counter interface {
	// 计数增加1
	Inc()
	// 返回滑动窗口内计数
	Count() int
}

// bucket 秒间隔记数
type bucket struct {
	timestamp int64
	amount    int
}

type counter struct {
	window   int // 窗口大小
	r        *ring.Ring
	increase chan int64
	mutex    sync.Mutex
}

// NewCounter 创建一个滑动窗口为window秒的计算器
func NewCounter(window int) Counter {
	c := &counter{window: window, increase: make(chan int64, window)}

	// 初始化ring.Ring
	r := ring.New(window + 1)
	now := time.Now().Unix()
	for i := 0; i < r.Len(); i++ {
		r.Value = &bucket{timestamp: now - int64(i)}
		r = r.Prev()
	}
	c.r = r.Next()

	go c.inc()

	return c
}

func (c *counter) Inc() {
	c.increase <- time.Now().Unix()
}

func (c *counter) Count() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now().Unix()

	r := c.r
	b := r.Value.(*bucket)
	if b.timestamp == now {
		r = r.Prev()
	}

	total := 0
	for i := 0; i < c.window; i++ {
		b = r.Value.(*bucket)
		if now-b.timestamp <= int64(c.window) {
			total += b.amount
		}
		r = r.Prev()
	}

	return total
}

func (c *counter) inc() {
	for {
		ts := <-c.increase
		r := c.r

		b := r.Value.(*bucket)
		if b.timestamp == ts {
			b.amount++
		} else {
			c.r = c.r.Next()
			r = c.r
			b = r.Value.(*bucket)
			b.timestamp = ts
			b.amount = 1
		}
	}
}
