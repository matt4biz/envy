package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	a, ok := os.LookupEnv("a")

	if ok {
		fmt.Println("a:", a)
	}

	fmt.Println("child waiting ...")

	select {
	case s := <-done:
		fmt.Println("child got", s)
	case <-time.After(10 * time.Second):
		fmt.Println("child timed out")
	}
}
