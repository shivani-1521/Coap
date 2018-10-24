package Coap

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

