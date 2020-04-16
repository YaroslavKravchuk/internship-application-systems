/*

Ping application created by Yaroslav Kravchuk
4/15/20

*/
package main
import (
  "flag"
  "os"
  "os/signal"
	"fmt"
  "strings"
	"log"
	"net"
	"time"

	"golang.org/x/net/icmp"
  "golang.org/x/net/ipv4"
  //"golang.org/x/net/ipv6"
)
func main() {

  // Get input
  flag.Parse()
  input := flag.Arg(0)

  // Set ip to input, regardless if input is IP or address
  ip := net.ParseIP(input)
  ipNet, err := net.ResolveIPAddr("ip", input)

  // Set ip to address IP if input is address
  if strings.ContainsAny(input, "abcdefghijklmnopqrstuvwxyz") {
    if err != nil {
      fmt.Println("Could not resolve address ", input)
      os.Exit(1)
    }
    ip = ipNet.IP
    fmt.Println(ip)
  } else {
    if input == "" {
      fmt.Println("No IP input received")
      os.Exit(1)
    }
  }

  // Keep track of packets and packet loss
  received := 0
  missed := 0
  numPings := 0

  // Create method to give staticstics and exit at ctrl-c
  sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	go func() {
    for _ = range sigchan {
      fmt.Printf("\n------%v Ping Statistics ------\n", input)
  	  fmt.Printf("%d packets transmitted, %d packets received, %d%% loss\n",
                  numPings, numPings, missed/(numPings*10)   )
  	  os.Exit(0)
    }
	}()

  // Loop ping requests
  for {
    conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
    if err != nil {
      os.Exit(1)
    }

    defer conn.Close()

    messICMP := icmp.Message{
      Type: ipv4.ICMPTypeEcho, Code: 0,
      Body: &icmp.Echo{
        ID: os.Getpid() & 0xffff, Seq:1,
        Data: []byte("aaaaaaaaaa"),
      },
    }

    message, err := messICMP.Marshal(nil)
    if err != nil {
      log.Fatal(err)
    }

    start := time.Now()

    writen, err := conn.WriteTo(message, &net.IPAddr{IP: ip})
  	if err != nil {
  		log.Fatalf("Error Writing, %s", err)
  	} else if writen != len(message) {
  		fmt.Println("Some packets were lost")
  	}

    readByte := make([]byte, 1028)
  	numCopied, packet, err := conn.ReadFrom(readByte)
    if err != nil {
  		log.Fatal(err)
  	}
    received += numCopied
    numPings++

    parsed, err := icmp.ParseMessage(1, readByte[:numCopied])
  	if err != nil {
  		log.Fatal(err)
  	}

    duration := time.Since(start)


    switch message := parsed.Body.(type) {
      case *icmp.Echo:
        missed += 10 - len(message.Data)
        for i := 0; i < len(message.Data); i++ {
          if message.Data[i] != 'a' {
            missed++
          }
        }
        fmt.Printf("\n%d bytes from %v: time=%d ms",
                      numCopied, packet, duration/1000000)

  	   default:
  		  log.Printf("error occured got %+v", parsed)
  	}

    time.Sleep(time.Second)
  }
}
