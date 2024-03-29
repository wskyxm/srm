package main

import (
	"fmt"
	"github.com/wskyxm/srm"
	"io"
	"net/http"
	"time"
)

type CustomData struct {
	Data1 string `json:"data_1"`
	Data2 string `json:"data_2"`
}

func callback(info *srm.SystemInfo) interface{} {
	data := CustomData{Data1: "11111", Data2: "22222"}
	return data
}

func report(w http.ResponseWriter, r *http.Request) {
	buf, _ := io.ReadAll(r.Body)
	fmt.Println(string(buf))
}

func main() {

	go func() {
		time.Sleep(time.Second * 3)

		srmobj := srm.NewSystemResourceMonitor(srm.SRMConfig{
			ListenAddr: ":10008",
			ReportAddr: "http://192.168.9.43:20008/test",
			ReportInterval: 10,
		})

		srmobj.Run()
	}()

	mux := http.ServeMux{}
	mux.HandleFunc("/test", report)
	http.ListenAndServe(":20008", &mux)
}
