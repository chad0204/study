package main

import (
	"fmt"
	"time"
)

type RateLimit struct {
	fillInterval time.Duration
	capacity     int
	bucket       chan struct{}
}

func NewRateLimit(fillInterval time.Duration, capacity int) *RateLimit {
	bucket := make(chan struct{}, capacity)
	go func() {
		ticker := time.NewTicker(fillInterval)
		for {
			select {
			case <-ticker.C:
				select {
				case bucket <- struct{}{}:
				default:
					//令牌桶满, 直接丢弃
				}
				//fmt.Printf("bucket token count: %v \n", len(bucket))
			}
		}
	}()
	return &RateLimit{
		capacity:     capacity,
		fillInterval: fillInterval,
		bucket:       bucket,
	}
}

func (limit *RateLimit) TakeAvailable(isBlock bool) bool {
	var result bool
	if isBlock {
		select {
		case <-limit.bucket:
			result = true
		}
	} else {
		select {
		case <-limit.bucket:
			result = true
		default:
			result = false
		}
	}
	return result
}

func main() {
	fillInterval := time.Millisecond * 10
	capacity := 100

	limit := NewRateLimit(fillInterval, capacity)

	for {
		fmt.Printf("get token: %v \n", limit.TakeAvailable(true))
	}

}
