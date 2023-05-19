package main

import (
	"fmt"
	"sync"
)

func processBusiness(id int, wg *sync.WaitGroup, stopChan chan bool) {
	defer wg.Done()

	// 模拟业务处理，假设某个线程处理业务返回true时立即停止其他线程
	// 这里使用一个简单的判断条件，当id为5时返回true
	if id == 5 {
		// 发送停止信号给其他线程
		stopChan <- true
		fmt.Printf("线程 %d 处理业务返回了 true\n", id)
		return
	}

	fmt.Printf("线程 %d 处理业务\n", id)
}

func main() {
	stopChan := make(chan bool) // 用于接收停止信号的通道
	var wg sync.WaitGroup

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go processBusiness(i, &wg, stopChan)
	}

	// 等待所有线程完成
	wg.Wait()

	// 关闭停止信号通道
	close(stopChan)
}
