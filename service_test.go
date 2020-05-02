package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestDelay(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Log(r.Intn(3))
	t.Log(r.Intn(3))
	t.Log(r.Intn(3))

}
