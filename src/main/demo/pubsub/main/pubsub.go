package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type (
	subscriber chan interface{}         //订阅者是一个管道
	topicFunc  func(v interface{}) bool //主题是一个过滤函数
)

//发布者对象
type Publisher struct {
	m           sync.RWMutex             //读写锁
	buffer      int                      //订阅队列得缓存大小
	timeout     time.Duration            //发布超时时间
	subscribers map[subscriber]topicFunc //订阅者信息
}

func NewPublisher(timeout time.Duration, buffer int) *Publisher {
	return &Publisher{
		timeout:     timeout,
		buffer:      buffer,
		subscribers: make(map[subscriber]topicFunc),
	}
}

// 订阅主题
func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	//创建订阅通道
	sub := make(chan interface{}, p.buffer)
	p.m.Lock()
	p.subscribers[sub] = topic
	p.m.Unlock()
	return sub
}

// 订阅所有主题
func (p *Publisher) SubscribeAll() chan interface{} {
	return p.SubscribeTopic(nil)
}

// 退出订阅
func (p *Publisher) Evict(sub chan interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	delete(p.subscribers, sub)
	close(sub)
}

// 关闭发布者, 即关闭所有订阅者
func (p *Publisher) Close() {
	p.m.RLock()
	defer p.m.RUnlock()

	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}

//发布主题
func (p *Publisher) Publish(v interface{}) {
	p.m.RLock()
	defer p.m.RUnlock()
	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		go p.sendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

// 发布主题（内部）, 设置超时时间
func (p *Publisher) sendTopic(
	sub subscriber,
	topic topicFunc,
	val interface{},
	wg *sync.WaitGroup,
) {
	if topic != nil && !topic(val) {
		return
	}
	//超时时间原理, select一直阻塞, 直到有case触发使select阻塞结束。
	//由于timeout后time会发送一个消息到chan让case2有效, 所以select至多会阻塞timeout时间, 在timeout内如果case1有效, 则select走case1。
	select {
	case sub <- val:
	case <-time.After(p.timeout):
	}
	defer wg.Done()
}

func main() {
	p := NewPublisher(1e9, 10)
	defer p.Close()

	//订阅所有主题的订阅者
	allSub := p.SubscribeAll()
	go func() {
		for msg := range allSub {
			fmt.Printf("recv msg from topic all. msg = %v \n", msg)
		}
	}()

	//订阅主题“golang”的订阅者
	golangSub := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			//s 包含 golang
			return strings.Contains(s, "golang")
		}
		return false
	})
	go func() {
		for msg := range golangSub {
			fmt.Printf("recv msg from topic golang. msg = %v \n", msg)
		}
	}()

	p.Publish("hello,  world!")
	p.Publish("hello, golang!")

	time.Sleep(3 * time.Second)
}
