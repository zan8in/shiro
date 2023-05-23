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

	result, err := s.RunMulti(shiro.Options{
		TargetFile:        "./1.txt",
		ShiroKeysFile:     "./keys.txt",
		RateLimitKey:      time.Duration(20),
		RateLimitTarget:   time.Duration(10),
		ConcurrencyKey:    2,
		ConcurrencyTarget: 6,
	})
	if err != nil {
		panic(err)
	}

	for v := range result {
		fmt.Printf("target: %s\n", v.Target)
		fmt.Printf("shiroKey: %s\n", v.ShiroKey)
		fmt.Printf("rememberMe: %s\n", v.RememberMe)
		fmt.Println("-----------------------")
	}

}
