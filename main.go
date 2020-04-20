package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

var ping = []string{"Ping for IPv4 addresses"}

func info() {
	app.Name = "Ping CLI"
	app.Usage = "A CLI for sending and recieving ICMP echo requests and replies"
	app.Author = "Kshitij Sharma"
	app.Version = "1.0.0"
}

func commands(){
	app.Commands = []cli.Command{
		{
			Name: "ping",
			Usage: "A CLI for sending and recieving ICMP requests to a hostname"
			Func:
		}
	}
}
func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
