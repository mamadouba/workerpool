package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"workerpool/poolworker"
	"workerpool/tasks"
)

func main1() {
	c := make(chan int, 2)
	c <- 1
	c <- 2
	a, b := <-c
	fmt.Println(a, b)
	close(c)
	a, b = <-c
	fmt.Println(a, b)
	a, b = <-c
	fmt.Println(a, b)

}
func main() {
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
