package cm

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
)

type RecvMsg struct {
	FileName   string `json:"file_name"`
	FileLength uint64 `json:"fileLength"`
	Offset     uint64 `json:"offset"`
	RemoteAddr string `json:"remoteAddr"`
}

func FileRecv(recv *RecvMsg) {
	conn, err := net.Dial("tcp", recv.RemoteAddr)
	defer conn.Close()
	if err != nil {
		log.Println("fileRecv err : connect err")
		return
	}
	var (
		file *os.File
		er   error
	)
	file, er = os.Open(recv.FileName)
	if er != nil {
		file, er = os.Create(recv.FileName)
	}
	defer file.Close()
	bfWt := bufio.NewWriter(file)
	for {
		data := make([]byte, 1024*128)
		n, err := conn.Read(data)
		if nil != err && err != io.EOF {
			log.Println(err.Error())
			return
		}
		if nil != err || n == 0 {
			bfWt.Flush()
			break
		}
		bfWt.Write(data[:n])
	}
}

/*func fileRescv(conn *net.Conn,ch chan string) {
	var fileMap map[string]string
	fileMsg,result := base.Msg.(string)
	if !result {
		log.Println("filePushErr : fileMsg err")
		return
	}
	json.Unmarshal([]byte(fileMsg),fileMap)

}*/
