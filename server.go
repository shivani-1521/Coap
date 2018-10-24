package CoapPS

import(
	"log"
	"net"
	"github.com/dustin/go-coap"
)

type chanMapStringList map[*net.UDPAddr][]string
type stringMapChanList map[string][]*net.UDPAddr

type CoapServer struct{
	capacity int
	msgIndex uint16
	clientMapTopics chanMapStringList
	topicMapClients stringMapChanList
}

func NewCoapPubSubServer(maxChannel int) *CoapServer{
	coapServe := new(CoapServer)

	coapServe.capacity = maxChannel
	coapServe.clientMapTopics = make(map[*net.UDPAddr][]string, maxChannel)
	coapServe.topicMapClients = make(map[string][]*net.UDPAddr, maxChannel)
	coapServe.msgIndex = GetIPv4Int16() + GetLocalRandomInt()
	log.Println("Init msgid : ", coapServe.msgIndex)

	return coapServe
}

func(c *CoapServer)genMsgId() uint16{

	c.msgIndex =c.msgIndex + 1

	return c.msgIndex
}

func(c *CoapServer)removeSub(topic string, client *net.UDPAddr){

	removeIndexT2C := -1
	if val, exists := c.topicMapClients[topic]; exists{
		for k, v := range val{
			if v == client{
				removeIndexT2C = k
			}
		}

		if removeIndexT2C != -1{
			sliceClients := c.topicMapClients[topic]
			if len(sliceClients) > 1 {
				c.topicMapClients[topic] = append(sliceClients[:removeIndexT2C], sliceClients[removeIndexT2C+1:]...)
			}else{
				delete(c.topicMapClients, topic)
			}
		}
	}

	removeIndexC2T := -1

	if val, exists := c.clientMapTopics[client]; exists{
		for k, v := range val {
			if v == topic {
				removeIndexC2T = k
			}
		}

		if removeIndexC2T != -1 {
			sliceTopics := c.clientMapTopics[client]
			if len(sliceTopics) > 1 {
				c.clientMapTopics[client] = append(sliceTopics[:removeIndexC2T], sliceTopics[removeIndexC2T+1:]...)
			}else{
				delete(c.clientMapTopics, client)
			}
		}
	}
}

func(c* CoapServer)addSub(topic string, client *net.UDPAddr){

	topicFound := false
	if val, exists := c.topicMapClients[topic]; exists{
		for _, val := val {
			if v == client{
				topicFound = true
			}
		}
	}

	if topicFound == false {
		c.topicMapClients[topic] = append(c.topicMapClients[topic], client)
	}

	clientFound := false
	if val, exists := c.clientMapTopics[client]; exists{
		for _, v := range val{
			if v == topic {
				clientFound = true
			}
		}
	}

	if clientFound == false{
		c.clientMapTopics[client] = append(c.clientMapTopics[client], topic)
	}

}

func (c *CoapServer)publish(con *net.UDPConn, topic string, msg string){
	if clients, exists := c.topicMapClients[topic]; !exists{
		return
	}else{
		for _, client := range clients{
			c.publishMsg(con, client, topic, msg)
			log.Println("topic ", topic, " Publish to ", client, " message", msg)

		}
	}

	log.Println("Publish done")
}

func(c *CoapServer) handleCoapMessage(con *net.UDPConn, add *net.UDPAddr, m *coap.Message) *coap.Message{

	var topic string
	if m.Path() != nil{
		topic = m.Path()[0]
	}

	cmd := ParseUint8ToString(m.Option(coap.ETag))
	log.Println("cmd ", cmd , " topic: ", topic, " message ", string(m.Payload))
	log.Println("code ",  m.Code, " option ", cmd)

	if cmd == "ADDSUB" {
		log.Println("Add subscription topic ", topic, " in client ", add)
		c.addSub(topic, add)
		c.responseOK(con, add, m)
	}else if cmd == "REMSUB" {
		log.Println("Remove subscription topic ", topic, " in client ", add)
		c.removeSub(topic, add)
		c.responseOK(con, add, m)
	}else if cmd == "PUB" {
		log.Println(con, topic, string(m.Payload))
		c.responseOK(con, add, m)
	}else if cmd == "HB" {
		log.Println("Got heart beat from ", add)
		c.responseOK(con, a, m)
	}

	for k, v := range c.topicMapClients {
		log.Println("Topic ", k, " sub by client ", v)
	}

	return nil

}

func(c *CoapServer)ListenAndServe(udpPort string){
	log.Fatal(coap.ListenAndServe("UDP ", udpPort, coap.FuncHandler(func(con *net.UDPConn, add *net.UDPAddr, m *coap.Message) *coap.Message{
		return c.handleCoapMessage(con, add, m)
		})))
}

func(c *CoapServer)responseOK(con *net.UDPConn, add *net.UDPAddr, m *coap.Message) {
	m2 := coap.Message{
		Type:	coap.Acknowledgement,
		Code:	coap.Content,
		MessageID:	m.MessageID,
		Payload:	m.Payload,

	}

	m2.SetOption(coap.ContentFormat, coap.TextPlain)
	m2.SetOption(coap.LocationPath, m.Path())

	err := coap.Transmit(conn, add, m2)
	if err != nil{
		log.Printf("Error in transmission %v ", err)
		return
	}

}

func(c *CoapServer)publishMsg(con *net.UDPConn, add *net.UDPAddr, topic string, msg string){
	m := coap.Message{
		Type:	coap.Confirmable,
		Code:	coap.Content,
		MessageID:	c.genMsgId(),
		Payload:	[]byte(msg),

	}

	m.SetOption(coap.ContentFormat, coap.TextPlain)
	m.SetOption(coap.LocationPath, topic)

	log.Printf("Transmitting %v msg %s", m, msg)
	err := coap.Transmit(con, add, m)
	if err != nil{
		log.Printf("transmission error %v", err)
		return
	}
}
