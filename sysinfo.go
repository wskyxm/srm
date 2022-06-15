package srm

import (
	"encoding/json"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"time"
)

type systemInfo struct {
	CpuUsage	int64 `json:"cpu_usage"`
	MemUsage	int64 `json:"mem_usage"`
	TotalMemory	int64 `json:"total_memory"`
	FreeMemory	int64 `json:"free_memory"`
	Timestamp	int64 `json:"timestamp"`
	Custom		interface{} `json:"custom"`
}

func update(callback func()interface{}) *systemInfo {
	// 获取内存和CPU信息
	cpuinfo, _ := cpu.Percent(time.Millisecond * 200, false)
	meminfo, _ := mem.VirtualMemory()

	// 获取系统资源占用情况
	sysinfo := &systemInfo{Timestamp: time.Now().Unix()}
	sysinfo.TotalMemory = int64(meminfo.Total / 1024.0 / 1024.0)
	sysinfo.FreeMemory = int64(meminfo.Free / 1024.0 / 1024.0)
	sysinfo.MemUsage = 100 - int64(float64(sysinfo.FreeMemory) / float64(sysinfo.TotalMemory) * 100.0)
	sysinfo.CpuUsage = int64(cpuinfo[0])

	// 调用回调函数
	if callback != nil {sysinfo.Custom = callback()}

	// 返回结果
	return sysinfo
}

func (s *systemInfo)tostring() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}
