package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	var transmitted int = 0
	var received int = 0
	var running = true // for infinite requests
	/* parse command line args */
	hostname := os.Args[1]

	/* just for show */
	fmt.Println("Hostname: ", hostname)

	/* creating channels to recieve signal notifications */
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	/* registers given channel, in this apps case fo CRTL+C */
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	/* executes blocking recieve for signals */
	go func() {
		<-sigs
		fmt.Println()
		fmt.Printf("--- %s ping statistics ---\n", hostname)
		running = false
		done <- true
	}()

	/* runs requests in a while loop until interrupted */
	for running {
		transmitted++
		ping(hostname)
		received++ // accounting for CRTL+C interrupt between the ping funciton
		time.Sleep(1 * time.Second)
	}

	/* summary data */
	fmt.Printf("%d packets transmitted, %d packets recieved, %f%% packet loss\n", transmitted, received, packetloss)
	<-done
}
