package main

import (
	"encoding/binary"
	"fmt"
	"net"

	log "gitlab.xinhulu.com/platform/GoPlatform/logger"
)

type Frame struct {
	frameType uint8
	channelId uint16
	size      uint32
	playload  []byte
	frameEnd  uint8
}

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:5672")
	if err != nil {
		_ = log.Errorf("%q", err)
		return
	}
	byteNum, err := conn.Write([]byte{'A', 'M', 'Q', 'P', 0, 0, 9, 1})
	if err != nil {

	}
	log.Info("write %d Byte to Server", byteNum)

	//var buf = make([]byte, 1000)
	//readNum, err := conn.Read(buf)
	//if err != nil {
	//
	//}
	//log.Info("read %d Byte From Server", readNum)
	//fmt.Println(buf[:readNum])

	var f Frame
	if err = binary.Read(conn, binary.BigEndian, &f.frameType); err != nil {
		return
	}
	if err = binary.Read(conn, binary.BigEndian, &f.channelId); err != nil {
		return
	}
	if err = binary.Read(conn, binary.BigEndian, &f.size); err != nil {
		return
	}

	//conn.Read(f.playload)

	fmt.Println(f)

}
