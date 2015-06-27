package common // import "moetang.info/go/common"

import (
	"fmt"
	"testing"
	"time"
)

import ()

func ExampleAA(t *testing.T) {
	c := time.Tick(1 * time.Second)
	for now := range c {
		fmt.Printf("%v\n", now)
	}

	ticker := time.NewTicker(1 * time.Second)
	for now := range ticker.C {
		fmt.Printf("%v\n", now)
	}
}

func TestBBB(t *testing.T) {

	ticker := time.NewTicker(1 * time.Second)

	ccc := make(chan bool, 1)

	go func() {
		for {
			tt := <-ticker.C
			fmt.Println(tt)
			fmt.Println(time.Now())
			ss := 0
			for i := 0; i < 5000000000; i++ {
				ss += i
			}
			fmt.Println(ss)
			fmt.Println(time.Now())
			fmt.Println("========")
		}
	}()

	for {
		ccc <- true
		time.Sleep(3 * time.Second)
	}
}
