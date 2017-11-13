package cm

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Base struct {
	Err int
	Msg interface{}
}
type IpMsg struct {
	Addr string
}
type WinMsg struct {
	IsReady bool
	Port    uint8
}
type File struct {
	Name  string
	Path  string
	Size  int64
	IsDir bool
}

/**
扫描端口监听
*/
func ScanPort(scanPort int) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+string(scanPort))
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return //终止程序
	}
	defer conn.Close()
	for {
		whenRecv(conn)
		log.Println("send Msg finished")
	}
}

/**
错误检查
*/
func checkError(err error) {
	if err != nil {
		log.Println("Error: %s", err.Error())
	}
}

/**
当收到扫描消息时的处理
*/
func whenRecv(conn *net.UDPConn) {
	baseMsg := new(Base)
	ipM := new(IpMsg)
	var buf [255]byte
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if nil != err {
		return
	}

	log.Println("received msg from :", addr)
	ipM.Addr = conn.LocalAddr().String()
	baseMsg.Err = 0
	baseMsg.Msg = ipM
	arr, erro := json.Marshal(baseMsg)
	if nil != erro {
		return
	}
	_, err = conn.WriteToUDP(arr, addr)
	checkError(err)
}