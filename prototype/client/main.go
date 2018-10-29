package main

import (
	"flag"
	"fmt"
	"log"
	. "github.com/shivani-1521/CoapPS"
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

	client := NewCoapPubsubClient("localhost:5683")
	if client == nil {
		log.Fatalln("Cannot connect to server")
	}

	if cmd == "ADDSUB" {
		ch, err := client.AddSub(topic)
		log.Println(" ch ", ch, " err=", err)
		log.Println("Got pub from topic ", topic, " pub ", <-ch)
	}
	log.Println("Done")
}
