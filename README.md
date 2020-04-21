# Basic Ping CLI Application

## Language Style & Design Decision Overview

1. The goal when building this project was to emphasize neat and organized code that is separated into blocks for not only easy reading, but debugging as well. I decided to challenge myself and try using Go for the first time, considering this as a mini personal project so I could keep myself busy and productive during quarantine. To my surprise, with previous python experience and alot of linux ext2 development in C from my operating systems class, the transition to using this language wasn't bad at all. In fact, some of the golang libraries were absolute saviours when making this app.

It truly amazed me that before I would be running terminal commands for fun, and now I was able to make something that models ping- thanks to this challenge (and go libraries).

2. My implementation of this project was split into two parts:
  - **ping.go** - This is the main file that does all the brainwork. Broken down into some steps, it starts listening for ICMP replies, sends a message, retrieves and parses the messages and feeds back important metadata such as loss, RTT times, and the IP endpoint address.
  - **main.go** - This file takes care of parsing the input through the terminal and running the pings in an infinite loop with an additional feature of the CRTL+C interrupt which when triggered, prints out some overall statistics modeled after the way we see ping in terminal.
  
    

## Software Dependencies
 I used to go1.12.2.darwin-amd64.pkg to install Go from golang.org/dl/ and did all the development on MacOS.

## To Run Application:

1. ```cd path/to/folder/``` (if downloaded as ping_icmp_cli just cd into the folder)
2. ```go get``` (might pop up with a folder warning or error but it should be fine)
3. ```go build``` (should have no error output)
4. ```sudo ./foldername hostname```

Note: hit ```CTRL+C``` to exit out of the loop!

## Example Input:
Some example tests you can run are:
- ```sudo ./ping_icmp_cli google.com```
- ```sudo ./ping_icmp_cli www.apple.com```
- ```sudo ./ping_icmp_cli 151.101.1.67```

![Example Outputt](https://ibb.co/hcTFcvJ)

## Future Goals 
Some things I would like to work and add on to this application would be support for IPv6 as it can only handle IPv4 as of right now. My thought process was to use a flag in the command line to signify whether an IPv6 or an IPv4 address is being passed in, so that I could dynamically set the ListenPacket, and ResolveIPAddress to handle IPv6 based on the input flag, however, I struggled with syntax issues and proper package use.

## Requirements
1. Use one of the specified languages
Please choose from among C/C++/Go/Rust. If you aren't familiar with these languages, you're not alone! Many engineers join Cloudflare without specific langauge experience. Please consult A Tour of Go or The Rust Programming Language.

2. Build a tool with a CLI interface
The tool should accept as a positional terminal argument a hostname or IP address.

3. Send ICMP "echo requests" in an infinite loop
As long as the program is running it should continue to emit requests with a periodic delay.

4. Report loss and RTT times for each message
Packet loss and latency should be reported as each message received.

## Thank you
Special shout out to CloudFlare for the unique brain buster, happy quarantine!