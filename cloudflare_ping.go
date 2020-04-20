/*
* cloudflare_ping.go
*
* A small Ping CLI application for MacOS or Linux
* This CLI app accepts a hostname or an IP address as its argument,
* then sends ICMP "echo requests" in a loop to the target while
* receiving "echo reply" messages
*
* Author: Kshitij Sharma
* Github: https://github.com/Kshitij-Sharma
*
 */

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	ProtocolICMP = 1 // Retrieved this constant from https://godoc.org/golang.org/x/net/internal/iana
)

/* ICMP is an internet layer protocol used by network devices to diagnose
network communication issues, used to determine whether or not data
is reaching its intended destination in a timely manner

The steps followed: (based on https://www.geeksforgeeks.org/ping-in-c/)
1.) Start listening for icmp replies
2.) Do a DNS lookup - converts website into IP address form
3.) In a loop, (send, wait, and recieve and ICMP per second)
4.) Print out data, finished
*/

/** listenInit
 * Intializes listening for ICMP replies
 *
 * Inputs: starts listening for icmp replies
 * Outputs: none
 * Side Effects: none
 * */

func listenInit(address string, c net.PacketConn, err error) {
	c, err = icmp.ListenPacket("ip4:icmp", address) // privileged ICMP endpoint, needs ip4:icmp
}

/** getHostByName
 * Intializes listening for ICMP replies
 *
 * Inputs: hostname (DNS or IP binary-format)
 * Outputs: none
 * Side Effects: none
 * */
func getHostByName(hostname string, err error, address *net.IPAddr) {
	address, err = net.ResolveIPAddr("ip4", hostname)
}

/** makeICMP
 * Makes the ICMP message to be sent
 *
 * Inputs: message to be populated, messageBinary buffer, error
 * Outputs: none
 * Side Effects: populats message, messagBinary
 * */
func makeICMP(message icmp.Message, messageBinary []byte, err error) {
	message = icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1, // the pid is a 32 bit so masked to grab proper data
			Data: []byte(""),
		},
	}
	messageBinary, err = message.Marshal(nil) // gives a binary encoding of the ICMP message
}

/** ping
 * Sends and retrieves ICMP messages per second
 *
 * Inputs: hostname (as a DNS or IP address)
 * Outputs: 0 on success, [null, 0, err] on failure
 * Side Effects: none
 * */
func ping(hostname string) (*net.IPAddr, time.Duration, float64, error) {
	var connection net.PacketConn // this confused me, but is basically a stream-oriented network conneeciton for packets
	var err error
	var ipv4Addr = "0.0.0.0"
	var myAddress *net.IPAddr
	var myMessage icmp.Message
	var myBinaryMessage []byte

	/* step 1: intialize listening for ICMP replies */
	listenInit(ipv4Addr, connection, err)
	if err != nil {
		return nil, 0, 0, err
	}
	//
	//
	// I only want to close WHEN CONTROL C is pressed
	//
	//

	defer connection.Close()

	/* step 2: do a DNS lookup on the hostname */
	getHostByName(hostname, err, myAddress)
	if err != nil {
		fmt.Printf("ping: cannot resolve %s: Unknown host")
		return nil, 0, 0, err
	}

	/* step 3: send an ICMP message and retrieve it every second */
	for true {
		/* make the ICMP message to send */
		makeICMP(myMessage, myBinaryMessage, err) // myBinaryMessage holds the binary encoding of ICMP message myMessage
		if err != nil {
			return myAddress, 0, 0, err
			break
		}

		/* send ICMP message */
		messageStart := time.Now()
		numBytesWritten, err := connection.WriteTo(myBinaryMessage, myAddress)

		if err != nil {
			return myAddress, 0, 0, err
			break
		} else if numBytesWritten != len(myBinaryMessage) {
			return myAddress, 0, 0, fmt.Errorf("got %v; want %v", numBytesWritten, len(myBinaryMessage))
			break
		}

		/* wait for ICMP message reply */
		replyBuffer := make([]byte, 1500)                                  // allocating space in buffer to read bytes into
		err = connection.SetReadDeadline(time.Now().Add(10 * time.Second)) // set a read time limit to 10 seconds

		if err != nil {
			return myAddress, 0, 0, err
			break
		}
		numBytesRead, peer, err := connection.ReadFrom(replyBuffer)
		if err != nil {
			return myAddress, 0, 0, err
		}
		messageDuration := time.Since(messageStart)

		/* step 4: print out data */
		retrievedMessage, err := icmp.ParseMessage(ProtocolICMP, replyBuffer[:numBytesRead])
		if err != nil {
			return myAddress, 0, 0, err
		}
		switch retrievedMessage.Type {
		case ipv4.ICMPTypeEchoReply:
			packetsRead := float64(numBytesRead)
			totalPackets := float64(len(myBinaryMessage))

			var packetloss float64
			packetloss = packetsRead / totalPackets

			if packetloss == 1 {
				packetloss = 0
			} else {
				packetloss = (1 - packetloss) * 100
			}
			fmt.Printf("Ping %s (%s): %s %v%% packet loss\n", hostname, myAddress, messageDuration, packetloss)
			return myAddress, messageDuration, packetloss, nil
		default:
			return myAddress, 0, 0, fmt.Errorf("got %+v from %v; want echo reply", retrievedMessage, peer)
		}
		time.Sleep(1 * time.Second)
	}
	return myAddress, 0, 0, err
}
