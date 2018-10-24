package main

import (
	"log"
	. "github.com/shivani-1521/Coap"
)

func main(){
	log.Println("Server start ...")
	serve := NewCoapPubSubServer(1024)
	serve.ListenAndServe(":5683")
}