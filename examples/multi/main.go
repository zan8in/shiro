package main

import (
	"fmt"

	"github.com/zan8in/shiro"
)

func main() {
	s, err := shiro.NewShiro()
	if err != nil {
		panic(err)
	}

	result, err := s.RunMulti(shiro.Options{
		TargetFile: "./1.txt",
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
