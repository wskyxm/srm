# srm
system resource monitor

# SRMConfig
```
type SRMConfig struct {
	ListenAddr		string              // 监听地址，被动查询当前系统资源占用信息
	ReportAddr		string              // 主动上报地址
	Callback		func()interface{}   // 回调函数，添加自定义的数据到系统信息中
	ReportInterval	int64               // 主动上报时间间隔，单位秒
}
```

# 主动上报
按设定的时间间隔POST系统资源占用信息

```
{
	"cpu_usage": 17, // CPU使用率
	"mem_usage": 48, // 内存使用率
	"total_memory": 16310, // 总内存，单位MB
	"free_memory": 8576, // 空闲内存，单位MB
	"timestamp": 1655262283, // 时间戳
	"custom": {
		"data_1": "11111", // 自定义数据
		"data_2": "22222" // 自定义数据
	}
}
```

# 查询接口
GET http://xxx.xxx.xxx.xxx/info

```
{
	"cpu_usage": 17, // CPU使用率
	"mem_usage": 48, // 内存使用率
	"total_memory": 16310, // 总内存，单位MB
	"free_memory": 8576, // 空闲内存，单位MB
	"timestamp": 1655262283, // 时间戳
	"custom": {
		"data_1": "11111", // 自定义数据
		"data_2": "22222" // 自定义数据
	}
}
```

# 调用示例
```
srm.Run(srm.SRMConfig{
	ListenAddr: ":10008",
	ReportAddr: "http://192.168.9.43:20008/test",
	ReportInterval: 10,
	Callback: callback,
})
```