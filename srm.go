package srm

import (
	"net/http"
	"strings"
	"time"
)

type SRMConfig struct {
	ListenAddr		string
	ReportAddr		string
	Callback		func()interface{}
	UpdateInterval	int64
	ReportInterval	int64
}

type systemResourceMonitor struct {
	res	*systemResource
	cfg SRMConfig
}

func Run(config SRMConfig) {
	// 参数检查
	if config.UpdateInterval <= 0 {config.UpdateInterval = 10}
	if config.ReportInterval <= 0 {config.ReportInterval = 10}
	if config.ListenAddr == "" {return}

	// 初始化监控对象
	srmobj := systemResourceMonitor{cfg: config}
	srmobj.res = NewSysMonitor(config.UpdateInterval, config.Callback)

	// 定时上报资源信息
	go srmobj.report()

	// 初始化HTTP路由
	mux := http.ServeMux{}
	mux.HandleFunc("/info", srmobj.info)

	// 启动HTTP服务
	http.ListenAndServe(config.ListenAddr, &mux)
}

func (s *systemResourceMonitor)info(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(s.res.get()))
}

func (s *systemResourceMonitor)report() {
	// 参数检查
	if s.cfg.ReportAddr == "" {return}

	// 上报资源信息
	doreport := func() {
		http.Post(s.cfg.ReportAddr, "application/json", strings.NewReader(s.res.get()))
		time.Sleep(time.Second * time.Duration(s.cfg.ReportInterval))
	}

	// 循环执行
	for {doreport()}
}