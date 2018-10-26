package CoapPS

import(
	
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

func init(){
	rand.Seed(time.Now().UnixNano())
}

func parseUint8ToString(intf interface{})string{
	val, ok := intf.([]uint8)
	if ok {
		return string(val)
	}else {
		return ""
	}
}

func GetLocalRandInt() uint16{
	rand.Seed(time.Now().UnixNano())
	return uint16(rand.Intn(1000))
}

func GetIPv4Int16() uint16 {

	ifaces, error := net.Interfaces()
	
	if error != nil {
		log.Println("No network ", error)
		return 0
	}

	for _, i := range ifaces {
		if strings.Contains(i.Name, "en0") {
			addrs, error := i.Addrs()
			// handle err
			if error != nil {
				log.Println("No IP:", error)
				return 0
			}

			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				
				if ip[0] == 0 {
					
					var myIP uint16
					myIP = uint16(ip[12])<<8 + uint16(ip[13])<<7 + uint16(ip[14])<<6 + uint16(ip[13])<<6
					return myIP
				}
			}
		}
	}
	return 0
}