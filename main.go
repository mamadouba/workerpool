package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"workerpool/datastruct"
	"workerpool/poolworker"
	"workerpool/tasks"
)

func main() {

	var stack = datastruct.Stack{}
	stack.Push(1)
	stack.Push(5)
	stack.Push(9)

	for !stack.IsEmpty() {
		fmt.Println(stack)
		stack.Pop()
	}

	var q = datastruct.Queue{}
	q.Enque(1)
	q.Enque(2)
	q.Enque(4)
	q.Enque(5)

	for !q.IsEmpty() {
		fmt.Println(q)
		q.Deque()
	}
}

func main1() {
	const maxTask = 30
	d := poolworker.New(4, 50).Start()
	for i := 1; i <= maxTask; i++ {
		task := poolworker.NewTask("fibo", tasks.Fibn, rand.Intn(50))
		d.Queue(task, time.Duration(2*time.Second))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("press CTRL+C to quit")
	<-quit
	d.Stop()
}
