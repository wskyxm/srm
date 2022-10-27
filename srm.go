package srm

import (
	"io"
	"net/http"
	"strings"
	"time"
)

type CollectInfoCallback func(*SystemInfo)
type ReportResultCallback func([]byte, error)

type SRMConfig struct {
	ListenAddr     string
	ReportAddr     string
	ReportInterval int64
}

type SystemResourceMonitor struct {
	config         SRMConfig
	OnCollectInfo  CollectInfoCallback
	OnReportResult ReportResultCallback
}

func NewSystemResourceMonitor(config SRMConfig) *SystemResourceMonitor {
	return &SystemResourceMonitor{config: config}
}

func (s *SystemResourceMonitor)Run() {
	// 参数检查
	if s.config.ReportInterval <= 0 {s.config.ReportInterval = 30}
	if s.config.ListenAddr == "" {return}

	// 定时上报资源信息
	go s.report()

	// 初始化HTTP路由
	mux := http.ServeMux{}
	mux.HandleFunc("/info", s.info)

	// 启动HTTP服务
	http.ListenAndServe(s.config.ListenAddr, &mux)
}

func (s *SystemResourceMonitor)info(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(update(s.OnCollectInfo).tostring()))
}

func (s *SystemResourceMonitor)report() {
	// 参数检查
	if s.config.ReportAddr == "" {
		return
	}

	// 上报资源信息
	doreport := func(data string) {
		resp, err := http.Post(s.config.ReportAddr, "application/json", strings.NewReader(data))
		if err != nil && s.OnReportResult != nil {s.OnReportResult(nil, err)}

		if err == nil && s.OnReportResult != nil {
			s.OnReportResult(io.ReadAll(resp.Body))
		}

		if err == nil {resp.Body.Close()}
		time.Sleep(time.Second * time.Duration(s.config.ReportInterval))
	}

	// 循环执行
	for {doreport(update(s.OnCollectInfo).tostring())}
}
