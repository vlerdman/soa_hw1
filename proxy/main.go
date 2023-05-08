package main

import (
	"log"
	"net"
	"os"
	"strconv"
)

var FormatToHost map[string]string

var GetResultPath = "get_result"

type Responce struct {
	Result string `json:"result"`
}

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("incorrect num of args provided: 2 is required")
	}
	port, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("incorrect port is provided")
	}
	
	FormatToHost = make(map[string]string)
    FormatToHost["json"] = "json_server:2001"
	FormatToHost["xml"] = "xml_server:2002"
	FormatToHost["yaml"] = "yaml_server:2003"
	
	ServerConn, _ := net.ListenUDP("udp", &net.UDPAddr{IP:[]byte{0,0,0,0},Port:port,Zone:""})
	defer ServerConn.Close()
	buf := make([]byte, 1024)



	for {
		n, addr, _ := ServerConn.ReadFromUDP(buf)

		log.Printf("receive bytes: %s", string(buf[0:n]))
		
		if n <= len(GetResultPath) || string(buf[0:len(GetResultPath)]) != GetResultPath {
			continue
		}
		format := string(buf[len(GetResultPath) + 1:n])
		host, ok := FormatToHost[format]
		if !ok {
			log.Printf("invalid format %s is provided", format)
			continue
		}

		udpServer, err := net.ResolveUDPAddr("udp", host)
        if err != nil {
			log.Printf("can't resolve udp server: %s", err)
			continue
		}


		conn, err := net.DialUDP("udp", nil, udpServer)
		if err != nil {
			log.Printf("can't connect to udp server: %s", err)
			continue
		}

		conn.Write([]byte("get_result"))

		received := make([]byte, 1024)
		conn.Read(received)

		ServerConn.WriteTo(received, addr)


		conn.Close()
	}
}
