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

	err = s.Run(shiro.Options{
		Target: "http://82.156.89.190:8090",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("shiro scanner finished")
}
