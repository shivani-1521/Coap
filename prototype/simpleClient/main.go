package main

import (
	"flag"
	"fmt"
	"log"
	"github.com/dustin/go-coap"
)

func main() {

	flag.Parse()
	if len(flag.Args()) < 3 {
		fmt.Println("Need more arguments")
		return
	}

	cmd := flag.Arg(0)
	topic := flag.Arg(1)
	msg := flag.Arg(2)

	fmt.Println(cmd, topic, msg)

	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: 12345,
		Payload:   []byte(msg),
	}

	req.SetOption(coap.ETag, cmd)
	req.SetOption(coap.MaxAge, 3)
	req.SetPathString(topic)

	c, err := coap.Dial("udp", "localhost:5683")
	if err != nil {
		log.Fatalf("Error in dialing: %v", err)
		return
	}

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error in sending request: %v", err)
	}

	if cmd == "PUB" {
		log.Println("Publish event send and leave ")
		return
	}

	log.Println("Pending for waiting publish result ")
	go func() {
		for {
			if rv != nil {
				if err != nil {
					log.Fatalf("Error receiving: %v", err)
				}
				log.Printf("Got %s", rv.Payload)
			}
			rv, err = c.Receive()
			log.Println("receiv:", rv, " err=", err)
		}
	}()

	for {
		
	}
log.Printf("Done\n")
}