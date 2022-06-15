package srm

import (
	"encoding/json"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"sync"
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

type systemResource struct {
	mutex    sync.RWMutex
	resinfo  systemInfo
	callback func()interface{}
}

func update(callback func()interface{}) systemInfo {
	// 获取内存和CPU信息
	cpuinfo, _ := cpu.Percent(time.Millisecond * 200, false)
	meminfo, _ := mem.VirtualMemory()

	// 获取系统资源占用情况
	sysinfo := systemInfo{Timestamp: time.Now().Unix()}
	sysinfo.TotalMemory = int64(meminfo.Total / 1024.0 / 1024.0)
	sysinfo.FreeMemory = int64(meminfo.Free / 1024.0 / 1024.0)
	sysinfo.MemUsage = 100 - int64(float64(sysinfo.FreeMemory) / float64(sysinfo.TotalMemory) * 100.0)
	sysinfo.CpuUsage = int64(cpuinfo[0])

	// 调用回调函数
	if callback != nil {sysinfo.Custom = callback()}

	// 返回结果
	return sysinfo
}

func NewSysMonitor(delay int64, f func()interface{}) *systemResource {
	// 保存回调函数
	s := &systemResource{callback: f, resinfo: update(f)}

	// 循环获取系统信息
	go s.update(delay)

	// 返回结果
	return s
}

func (s *systemInfo)tostring() string {
	buf, _ := json.Marshal(s)
	return string(buf)
}

func (s *systemResource)get() string {
	// 加锁
	s.mutex.Lock(); defer s.mutex.Unlock()

	// 返回结果
	return s.resinfo.tostring()
}

func (s *systemResource)update(delay int64) {
	for {
		// 获取系统信息
		s.mutex.Lock()
		s.resinfo = update(s.callback)
		s.mutex.Unlock()

		// 等待指定的时间
		time.Sleep(time.Second * time.Duration(delay))
	}
}
