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
	ReportInterval	int64
}

type systemResourceMonitor struct {
	cfg SRMConfig
}

func Run(config SRMConfig) {
	// 参数检查
	if config.ReportInterval <= 0 {config.ReportInterval = 30}
	if config.ListenAddr == "" {return}

	// 初始化监控对象
	srmobj := systemResourceMonitor{cfg: config}

	// 定时上报资源信息
	go srmobj.report()

	// 初始化HTTP路由
	mux := http.ServeMux{}
	mux.HandleFunc("/info", srmobj.info)

	// 启动HTTP服务
	http.ListenAndServe(config.ListenAddr, &mux)
}

func (s *systemResourceMonitor)info(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(update(s.cfg.Callback).tostring()))
}

func (s *systemResourceMonitor)report() {
	// 参数检查
	if s.cfg.ReportAddr == "" {return}

	// 上报资源信息
	doreport := func(data string) {
		http.Post(s.cfg.ReportAddr, "application/json", strings.NewReader(data))
		time.Sleep(time.Second * time.Duration(s.cfg.ReportInterval))
	}

	// 循环执行
	for {doreport(update(s.cfg.Callback).tostring())}
}