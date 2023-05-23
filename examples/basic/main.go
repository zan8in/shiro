package main

import (
	"fmt"
	"time"

	"github.com/zan8in/shiro"
)

func main() {
	s, err := shiro.NewShiro()
	if err != nil {
		panic(err)
	}

	result, err := s.Run(shiro.Options{
		Target:            "http://82.156.89.190:8090",
		RateLimitKey:      time.Duration(20),
		RateLimitTarget:   time.Duration(10),
		ConcurrencyKey:    2,
		ConcurrencyTarget: 6,
	})
	if err != nil {
		panic(err)
	}

	if result == nil {
		fmt.Println("result is nil")
		return
	}

	fmt.Printf("target: %s\n", result.Target)
	fmt.Printf("shiroKey: %s\n", result.ShiroKey)
	fmt.Printf("rememberMe: %s\n", result.RememberMe)

}
