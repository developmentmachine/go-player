package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/playing", playing)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalln(err)
	}
}

func playing(resp http.ResponseWriter, req *http.Request) {
	file, err := os.Open("d:\\test.mp4")
	if err != nil {
		resp.WriteHeader(500)
		return
	}
	stat, err := file.Stat()
	if err != nil {
		resp.WriteHeader(500)
		return
	}
	var minRange = "0"
	var maxRange = strconv.Itoa(int(stat.Size() - 1))
	if req.Header.Get("Range") != "" && strings.Contains(req.Header.Get("Range"), "-") {
		r := strings.Split(strings.ReplaceAll(req.Header.Get("Range"), "bytes=", ""), "-")
		minRange = r[0]
		if len(r) == 2 && r[1] != "" {
			maxRange = r[1]
		}
		resp.Header().Add("Content-Range", strings.Join([]string{"bytes ", r[0], "-", maxRange, "/", strconv.Itoa(int(stat.Size()))}, ""))
	}
	resp.Header().Add("Content-Length", strconv.Itoa(int(stat.Size())))
	resp.Header().Add("Content-Type", "video/mp4")
	resp.WriteHeader(206)
	seek, _ := strconv.ParseInt(minRange, 10, 64)
	file.Seek(seek, io.SeekStart)
	buffer := make([]byte, 204800)
	for {
		if read, err := file.Read(buffer); read != 0 && err != io.EOF {
			resp.Write(buffer)
		}
	}
}
