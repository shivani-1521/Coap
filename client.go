package CoapPubsub

import (
	"errors"
	"log"
	"time"
	"github.com/dustin/go-coap"
)

type subscripConnection struct {

	channel   chan string
	clientCon *coap.Conn

}

type CoapPubsubClient struct {

	msgIndex uint16
	serverAddress  string
	subList  map[string]subscripConnection

}

func NewCoapPubsubClient(servAddr string) *CoapPubsubClient {
	c := new(CoapPubsubClient)
	c.subList = make(map[string]subscripConnection, 0)
	c.serverAddress = servAddr

	c.msgIndex = GetIPv4Int16() + GetLocalRandomInt()
	log.Println("Init messageID ", c.msgIndex)
	go c.heartBeat()
	return c
}

func (c *CoapPubsubClient) AddSub(topic string) (chan string, error) {
	if val, exist := c.subList[topic]; exist {
		return val.channel, nil
	}

	conn, err := c.sendPubsubReq("ADDSUB", topic)
	if err != nil {
		return nil, err
	}

	subChan := make(chan string)
	go c.waitSubResponse(conn, subChan, topic)

	clientConn := subscripConnection{channel: subChan, clientCon: conn}
	c.subList[topic] = clientConn
	return subChan, nil
}

func (c *CoapPubsubClient) RemoveSub(topic string) error {
	if _, exist := c.subList[topic]; !exist {
		return nil
	}

	_, err := c.sendPubsubReq("REMSUB", topic)
	return err
}

func (c *CoapPubsubClient) sendPubsubReq(cmd string, topic string) (*coap.Conn, error) {
	Req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: c.getMsgID(),
		Payload:   []byte(""),
	}

	Req.SetOption(coap.ETag, cmd)
	Req.SetPathString(topic)

	conn, err := coap.Dial("udp", c.serverAddress)
	if err != nil {
		log.Printf(cmd, ">>Error dialing: %v \n", err)
		return nil, errors.New("Dial failed")
	}
	conn.Send(Req)
	return conn, err
}

func (c *CoapPubsubClient) waitSubResponse(conn *coap.Conn, ch chan string, topic string) {
	var rv *coap.Message
	var err error
	var keepLoop bool
	keepLoop = true
	for keepLoop {
		if rv != nil {
			if err != nil {
				log.Fatalf("Error receiving: %v", err)
			}
			log.Printf("Got %s", rv.Payload)
		}
		rv, err = conn.Receive()

		if err == nil {
			ch <- string(rv.Payload)
		}

		time.Sleep(time.Second)
		if _, exist := c.subList[topic]; !exist {
			log.Println("Loop topic:", topic, " already remove leave loop")
			keepLoop = false
		}
	}
}

func (c *CoapPubsubClient) getMsgID() uint16 {
	c.msgIndex = c.msgIndex + 1
	return c.msgIndex
}

func (c *CoapPubsubClient) heartBeat() {
	log.Println("Starting heart beat loop call")
	hbReq := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: c.getMsgID(),
		Payload:   []byte("Heart beat msg."),
	}

	hbReq.SetOption(coap.ETag, "HB")

	for {

		for k, conn := range c.subList {
			conn.clientCon.Send(hbReq)
			log.Println("Send the heart beat in topic ", k)
		}

		time.Sleep(time.Minute)
	}
}