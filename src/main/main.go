package main

import (
	"bytes"
	"cm"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

/**
开启服务
 */
func startServ()  {
	httpPort := 8888
	go cm.ScanPort(22555, httpPort)
	http.HandleFunc("/test", test)
	http.HandleFunc("/fileIO", fileIO)
	http.HandleFunc("/recv", fileRecv)
	http.HandleFunc("/list", listFile)
	err := http.ListenAndServe(":"+strconv.Itoa(httpPort), nil)
	if err != nil {
		log.Println("ListenAndServe: ", err.Error())
		return
	}
}

/**
列表文件信息
*/
func listFile(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	req.ParseForm()
	filePath := req.Form.Get("filePath")
	if filePath == "" {
		fList := GetLogicalDrives()
		var files [20]cm.FileMsg
		for i, v := range fList {
			files[i] = cm.FileMsg{
				Name:  v,
				Path:  v+"\\",
				Size:  0,
				IsDir: true,
			}
		}
		genList := cm.Base{
			Err: 0,
			Msg: files[:len(fList)],
		}
		results, _ := json.Marshal(genList)
		w.Write(results)
		return
	}
	f, _ := os.Open(filePath)
	fInfo, _ := f.Stat()
	if fInfo.IsDir() {
		osFList, _ := ioutil.ReadDir(filePath)
		var files [1024]cm.FileMsg
		for i, v := range osFList {
			bf := bytes.Buffer{}
			bf.WriteString(filePath)
			bf.WriteRune(os.PathSeparator)
			bf.WriteString(v.Name())
			fl,_ :=filepath.Abs(bf.String())
			files[i] = cm.FileMsg{
				Name:  v.Name(),
				Path:  fl,
				Size:  v.Size(),
				IsDir: v.IsDir(),
			}
		}
		genList := cm.Base{
			Err: 0,
			Msg: files[:len(osFList)],
		}
		results, _ := json.Marshal(genList)
		w.Write(results)
		return
	}
}

/**
接收文件
*/
func fileRecv(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	e := req.ParseForm()
	checkErr(e)
	var rm = new(cm.RecvMsg)
	err := json.Unmarshal([]byte(req.Form.Get("RecvMsg")), rm)
	checkErr(err)
	go cm.FileRecv(rm)
	results, _ := json.Marshal(cm.Base{
		Err: 0,
		Msg: "Is connected",
	})
	w.Write(results)
}

/**
连接测试
*/
func test(w http.ResponseWriter, req *http.Request) {
	results, _ := json.Marshal(cm.Base{
		Err: 0,
		Msg: "test finish",
	})
	w.Write(results)
}

func main() {
	startServ()
}

/**
文件流
 */
func fileIO(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	filePath := req.Form.Get("filePath")
	http.ServeFile(w,req,filePath)
}

/**
获取系统盘符
*/
func GetLogicalDrives() []string {
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	GetLogicalDrives := kernel32.MustFindProc("GetLogicalDrives")
	n, _, _ := GetLogicalDrives.Call()
	s := strconv.FormatInt(int64(n), 2)
	var drivesAll = []string{"A:", "B:", "C:", "D:", "E:", "F:", "G:", "H:", "I:", "J:", "K:", "L:", "M:", "N:", "O:", "P：", "Q：", "R：", "S：", "T：", "U：", "V：", "W：", "X：", "Y：", "Z："}
	temp := drivesAll[0:len(s)]
	var d []string
	for i, v := range s {
		if v == 49 {
			l := len(s) - i - 1
			d = append(d, temp[l])
		}
	}

	var drives []string
	for i, v := range d {
		drives = append(drives[i:], append([]string{v}, drives[:i]...)...)
	}
	return drives

}

func checkErr(err error) {
	if nil != err {
		return
	}
}
