package cm

import (
	"encoding/json"
	"log"
	"net"
	"strconv"
)

type Base struct {
	Err int         `json:"err"`
	Msg interface{} `json:"msg"`
}
type IpMsg struct {
	Port int64 `json:"port"`
}
type FileMsg struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"isDir"`
}

/**
扫描端口监听
*/
func ScanPort(scanPort int, httpPort int) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(scanPort))
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Println("Error listening", err.Error())
		return //终止程序
	}
	defer conn.Close()
	for {
		whenRecv(conn, httpPort)
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
func whenRecv(conn *net.UDPConn, httpPort int) {
	baseMsg := new(Base)
	ipM := new(IpMsg)
	var buf [255]byte
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if nil != err {
		return
	}
	log.Println("received msg from :", addr)
	ipM.Port = int64(httpPort)
	baseMsg.Err = 0
	baseMsg.Msg = ipM
	arr, erro := json.Marshal(baseMsg)
	if nil != erro {
		baseMsg.Err = 1
		return
	}
	udpAddr, err := net.ResolveUDPAddr("udp", (net.IP)(addr.IP).String()+":"+strconv.Itoa(22455))
	_, err = conn.WriteToUDP(arr, udpAddr)
	checkError(err)
}
