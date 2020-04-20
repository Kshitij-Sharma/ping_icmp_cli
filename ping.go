/*
* ping.go
*
* A small Ping CLI application for MacOS or Linux
* This CLI app accepts a hostname or an IP address as its argument,
* then sends ICMP "echo requests" in a loop to the target while
* receiving "echo reply" messages
*
* Note: Uses main.go to make an infinite loop of pings
*
* Author: Kshitij Sharma
* Github: https://github.com/Kshitij-Sharma/ping_icmp_cli
*
 */

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	ProtocolICMP = 1 // Retrieved this constant from https://godoc.org/golang.org/x/net/internal/iana
)

/* global variables */
var packetloss float32

/* ICMP is an internet layer protocol used by network devices to diagnose
network communication issues, used to determine whether or not data
is reaching its intended destination in a timely manner

The steps followed: (based on https://www.geeksforgeeks.org/ping-in-c/)
1.) Start listening for icmp replies
2.) Do a DNS lookup - converts website into IP address form
3.) In a loop, (send, wait, and recieve and ICMP per second)
4.) Print out data, finished
*/

/** makeICMP
 * Makes the ICMP message to be sent
 *
 * Inputs: message to be populated, messageBinary buffer, error
 * Outputs: none
 * Side Effects: populats message, messagBinary
 * */
func makeICMP() ([]byte, error) {
	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1, // the pid is a 32 bit so masked to grab proper data
			Data: []byte(""),
		},
	}
	messageBinary, err := message.Marshal(nil) // gives a binary encoding of the ICMP message
	return messageBinary, err
}

/** fatalErrhandler
 * Prints out fatal error and exits function if error exits
 *
 * Inputs: erorr
 * Outputs: none
 * Side Effects: none
 * */
func fatalErrHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/** ping
 * Handles ping requests by sending and recieving ICMP messages
 *
 * Inputs: hostname (as a DNS or IP address)
 * Outputs: 0 on success, [null, 0, err] on failure
 * Side Effects: none
 * */
func ping(hostname string) (*net.IPAddr, time.Duration, float64, error) {
	// var err error
	var ipv4Addr = "0.0.0.0"
	var myAddress *net.IPAddr
	var myBinaryMessage []byte

	/* step 1: intialize listening for ICMP replies */
	connection, err := icmp.ListenPacket("ip4:icmp", ipv4Addr)
	fatalErrHandler(err)

	defer connection.Close()

	/* step 2: do a DNS lookup on the hostname */
	myAddress, err = net.ResolveIPAddr("ip4", hostname)
	fatalErrHandler(err)

	/* step 3: send an ICMP message and retrieve */

	myBinaryMessage, err = makeICMP() // myBinaryMessage holds the binary encoding of ICMP message myMessage
	fatalErrHandler(err)

	/* send ICMP message */
	messageStart := time.Now()
	numBytesWritten, err := connection.WriteTo(myBinaryMessage, myAddress)

	fatalErrHandler(err)
	if numBytesWritten != len(myBinaryMessage) {
		err = fmt.Errorf("got %v; want %v", numBytesWritten, len(myBinaryMessage))
		fmt.Println(err.Error())
		return myAddress, 0, 0, err
	}

	/* wait for ICMP message reply */
	replyBuffer := make([]byte, 1500)                                  // allocating space in buffer to read bytes into
	err = connection.SetReadDeadline(time.Now().Add(10 * time.Second)) // set a read time limit to 10 seconds

	fatalErrHandler(err)
	// Retrieve and read ICMP Message
	numBytesRead, peer, err := connection.ReadFrom(replyBuffer)
	fatalErrHandler(err)
	messageDuration := time.Since(messageStart)

	/* step 4: print out data */
	retrievedMessage, err := icmp.ParseMessage(ProtocolICMP, replyBuffer[:numBytesRead])
	fatalErrHandler(err)

	switch retrievedMessage.Type {
	case ipv4.ICMPTypeEchoReply:
		packetsRead := float32(numBytesRead)
		totalPackets := float32(len(myBinaryMessage))

		packetloss = packetsRead / totalPackets

		// calculator for packet loss in terms of bytes, packetloss will hold a percentage value
		if packetloss == 1 {
			packetloss = 0
		} else {
			packetloss = (1 - packetloss) * 100
		}
		fmt.Printf("Ping from %s (%s): time=%s %v%% packet loss (bytes) \n", hostname, myAddress, messageDuration, packetloss)
	default:
		err = fmt.Errorf("got %+v from %v; want echo reply", retrievedMessage, peer)
		fmt.Println(err.Error())
	}
	return myAddress, 0, 0, err
}
