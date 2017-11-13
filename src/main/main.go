package main

import (
	"bufio"
	"cm"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
)

/**
列表文件信息
*/
func listFile(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	req.ParseForm()
	filePath := req.Form.Get("filePath")
	if filePath == "" {
		fList := GetLogicalDrives()
		var files [20]cm.File
		for i, v := range fList {
			files[i] = cm.File{
				Name:  v,
				Size:  0,
				IsDir: false,
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
		var files [1024]cm.File
		for i, v := range osFList {
			files[i] = cm.File{
				Name:  v.Name(),
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
get参数 : filePath 文件路径
下载指定文件
*/
func fileDownLoad(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	filePath := req.Form.Get("filePath")
	file, err := os.Open(filePath)
	log.Println(filePath)
	defer file.Close()
	if nil != err {
		fmt.Fprintf(w, err.Error())
		return
	}
	fileInfo, errs := file.Stat()
	if nil != errs {
		fmt.Fprintf(w, errs.Error())
		return
	}
	var offset int64 = 0
	var fileLen = fileInfo.Size()
	rag := req.Header.Get("Range")
	index := strings.Index(rag, "-")
	log.Println(index)
	switch index {
	case 0:
		r, _ := strconv.ParseInt(rag, 10, 64)
		offset = fileLen + r
	case len(rag) - 1:
		r, _ := strconv.ParseInt(strings.Replace(rag, "-", "", -1),
			10, 64)
		offset = r
	case -1:
	default:
		strSlice := strings.Split(rag, "-")[1:]
		offset, _ = strconv.ParseInt(strSlice[0], 10, 64)
		fileLen, _ = strconv.ParseInt(strSlice[1], 10, 64)
	}
	w.Header().Add("Accept-Ranges", "bytes")
	w.Header().Add("Content-disposition", "attachment;filename="+fileInfo.Name())
	w.Header().Add("Content-Length", strconv.FormatInt(fileLen-offset, 10))
	file.Seek(offset, 0)
	bfRd := bufio.NewReader(file)
	for {
		data := make([]byte, 1024*128)
		n, err := bfRd.Read(data)
		if nil != err && err != io.EOF {
			log.Println(err.Error())
			return
		}
		if nil != err || n == 0 {
			break
		}
		_, errw := w.Write(data[:n])
		if nil != errw {
			log.Println(errw.Error())
		}
	}
}

func main() {
	go cm.ScanPort(22555)
	http.HandleFunc("/list", listFile)
	http.HandleFunc("/fileDL", fileDownLoad)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Println("ListenAndServe: ", err.Error())
		return
	}
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
