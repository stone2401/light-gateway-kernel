package pcore

import (
	"sync/atomic"
	"time"
)

type Counter struct {
	count atomic.Int64
	cache chan Entry
}

type Entry struct {
	Count int64
	Time  time.Time
}

func NewCounter(cacheSize int) *Counter {
	c := &Counter{
		count: atomic.Int64{},
		cache: make(chan Entry, cacheSize),
	}
	go func() {
		// 等到５的整数时刻，例如　０，５，１０
		time.Sleep(time.Duration(5-time.Now().Second()%5) * time.Second)
		// 每５ｓ钟，将统计值写入缓存
		t := time.NewTicker(5 * time.Second)
		for {
			<-t.C
			entry := Entry{
				Count: c.count.Swap(0),
				Time:  time.Now(),
			}
			select {
			case c.cache <- entry:
			default:
				// 缓存已满，丢弃最早的一次数据
				<-c.cache
				c.cache <- entry
			}
		}
	}()
	return c
}

func (c *Counter) Inc() {
	c.count.Add(1)
}

// 获取统计
func (c *Counter) Gain() chan Entry {
	return c.cache
}

func (c *Counter) CounterHandler(ctx *Context) {
	ctx.Next()
	c.Inc()
}
