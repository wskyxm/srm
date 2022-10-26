package srm

import (
	"net/http"
	"strings"
	"time"
)

type SRMCallback func(*SystemInfo) interface{}

type SRMConfig struct {
	ListenAddr     string
	ReportAddr     string
	ReportInterval int64
}

type srmconfig struct {
	ListenAddr     string
	ReportAddr     string
	Callback       SRMCallback
	ReportInterval int64
}

type systemResourceMonitor struct {
	cfg srmconfig
}

func Run(config SRMConfig, callback SRMCallback) {
	// 参数检查
	if config.ReportInterval <= 0 {
		config.ReportInterval = 30
	}
	if config.ListenAddr == "" {
		return
	}

	// 初始化监控对象
	srmobj := systemResourceMonitor{cfg: srmconfig{
		ListenAddr:     config.ListenAddr,
		ReportAddr:     config.ReportAddr,
		ReportInterval: config.ReportInterval,
		Callback:       callback,
	}}

	// 定时上报资源信息
	go srmobj.report()

	// 初始化HTTP路由
	mux := http.ServeMux{}
	mux.HandleFunc("/info", srmobj.info)

	// 启动HTTP服务
	http.ListenAndServe(config.ListenAddr, &mux)
}

func (s *systemResourceMonitor) info(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(update(s.cfg.Callback).tostring()))
}

func (s *systemResourceMonitor) report() {
	// 参数检查
	if s.cfg.ReportAddr == "" {
		return
	}

	// 上报资源信息
	doreport := func(data string) {
		http.Post(s.cfg.ReportAddr, "application/json", strings.NewReader(data))
		time.Sleep(time.Second * time.Duration(s.cfg.ReportInterval))
	}

	// 循环执行
	for {
		doreport(update(s.cfg.Callback).tostring())
	}
}
