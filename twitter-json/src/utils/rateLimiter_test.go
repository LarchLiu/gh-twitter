package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestLimit(t *testing.T) {
	// go test -timeout 120s -run ^TestLimit$ ./src/utils -v
	limiter := NewRateLimiter(time.Minute, 20)
	// var limiter *RateLimiter = nil
	start := time.Now()

	for i := 0; i < 30; i++ {
		if limiter != nil {
			limiter.Limit()
		}
		fmt.Printf("i is %d \n", i)
		fmt.Println("test")
	}

	end := time.Now()

	d := end.Sub(start)
	fmt.Println("spends seconds: ", d.Seconds())
}
