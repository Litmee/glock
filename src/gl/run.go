package gl

import (
	"log"
	"net"
)

func Start() {
	listen, err := net.Listen("tcp", ":9098")
	if err != nil {
		panic(err)
	}
	log.Println("glock is start")
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("accept failed, err:%V\n", err)
			continue
		}
		go tcpHandler(&conn)
	}
}
