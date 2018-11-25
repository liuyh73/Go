package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/liuyh73/Go/udp"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	serverAddr := LFTP.SERVER_IP + ":" + strconv.Itoa(LFTP.SERVER_PORT)
	conn, err := net.Dial("udp", serverAddr)
	checkErr(err)

	defer conn.Close()
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()

		lineLen := len(line)

		n := 0
		for written := 0; written < lineLen; written += n {
			var toWrite string
			if lineLen-written > LFTP.SERVER_RECV_LEN {
				toWrite = line[written : written+LFTP.SERVER_RECV_LEN]
			} else {
				toWrite = line[written:]
			}
			n, err = conn.Write([]byte(toWrite))
			checkErr(err)

			fmt.Println("Write: ", toWrite)

			msg := make([]byte, LFTP.SERVER_RECV_LEN)
			n, err = conn.Read(msg)

			fmt.Println("Response: ", string(msg))
		}
	}
}
