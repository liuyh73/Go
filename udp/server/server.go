package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/liuyh73/LFTP"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	serverAddr := LFTP.SERVER_IP + ":" + strconv.Itoa(LFTP.SERVER_PORT)
	serverUDPAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	checkErr(err)

	conn, err := net.ListenUDP("udp", serverUDPAddr)
	checkErr(err)

	defer conn.Close()

	for {
		data := make([]byte, LFTP.SERVER_RECV_LEN)
		_, clientUDPAddr, err := conn.ReadFromUDP(data)

		if err != nil {
			fmt.Println(err)
			continue
		}
		strData := string(data)
		fmt.Println("Received: ", strData)

		upper := strings.ToUpper(strData)
		_, err = conn.WriteToUDP([]byte(upper), clientUDPAddr)

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Send: ", upper)
	}
}
