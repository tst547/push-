package cm

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Base struct {
	Err int         `json:"err"`
	Msg interface{} `json:"msg"`
}
type IpMsg struct {
	Addr string `json:"addr"`
	Port int64  `json:"port"`
}
type File struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"isDir"`
}

/**
扫描端口监听
*/
func ScanPort(scanPort int,httpPort int) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(scanPort))
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return //终止程序
	}
	defer conn.Close()
	for {
		whenRecv(conn,httpPort)
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
func whenRecv(conn *net.UDPConn,httpPort int) {
	baseMsg := new(Base)
	ipM := new(IpMsg)
	var buf [255]byte
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if nil != err {
		return
	}

	log.Println("received msg from :", addr)
	remote := addr.IP.String()
	ipM.Addr = localIp(remote[:strings.LastIndex(remote,".")+1])
	ipM.Port = int64(httpPort)
	baseMsg.Err = 0
	baseMsg.Msg = ipM
	arr, erro := json.Marshal(baseMsg)
	if nil != erro {
		return
	}
	_, err = conn.WriteToUDP(arr, addr)
	checkError(err)
}
func localIp(remote string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	i := 0
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if strings.Contains(ipnet.IP.String(),remote){
					return ipnet.IP.String()
				}
				// 检查ip地址判断是否回环地址
				i++
			}

		}
	}
	return ""
}