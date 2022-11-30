package main

import (
	"study/src/main/demo/sync/main/memotest"
	"testing"
)

var HTTPGetBody = memotest.HTTPGetBody

func Test(t *testing.T) {
	m := New(HTTPGetBody)
	memotest.Sequential(t, m)
}

func TestConcurrent(t *testing.T) {
	m := New(HTTPGetBody)
	memotest.Concurrent(t, m)
}
