package controller

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/tiger1103/gfast/v3/api/v1/system"
	"github.com/tiger1103/gfast/v3/library/libUtils"
	"os"
	"runtime"
	"strconv"
	"time"
)

var Monitor = sysMonitorController{
	startTime: gtime.Now(),
}

type sysMonitorController struct {
	BaseController
	startTime *gtime.Time
}

func (c *sysMonitorController) List(ctx context.Context, req *system.MonitorSearchReq) (res *system.MonitorSearchRes, err error) {
	cpuNum := runtime.NumCPU() //核心数
	var cpuUsed float64 = 0    //用户使用率
	var cpuAvg5 float64 = 0    //CPU负载5
	var cpuAvg15 float64 = 0   //当前空闲率

	cpuInfo, err := cpu.Percent(time.Duration(time.Second), false)
	if err == nil {
		cpuUsed, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", cpuInfo[0]), 64)
	}

	loadInfo, err := load.Avg()
	if err == nil {
		cpuAvg5, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", loadInfo.Load5), 64)
		cpuAvg15, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", loadInfo.Load5), 64)
	}

	var memTotal uint64 = 0  //总内存
	var memUsed uint64 = 0   //总内存  := 0 //已用内存
	var memFree uint64 = 0   //剩余内存
	var memUsage float64 = 0 //使用率

	v, err := mem.VirtualMemory()
	if err == nil {
		memTotal = v.Total
		memUsed = v.Used
		memFree = memTotal - memUsed
		memUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", v.UsedPercent), 64)
	}

	var goTotal uint64 = 0  //go分配的总内存数
	var goUsed uint64 = 0   //go使用的内存数
	var goFree uint64 = 0   //go剩余的内存数
	var goUsage float64 = 0 //使用率

	var gomem runtime.MemStats
	runtime.ReadMemStats(&gomem)
	goUsed = gomem.Sys
	goUsage = gconv.Float64(fmt.Sprintf("%.2f", gconv.Float64(goUsed)/gconv.Float64(memTotal)*100))
	sysComputerIp := "" //服务器IP

	ip, err := libUtils.GetLocalIP()
	if err == nil {
		sysComputerIp = ip
	}

	sysComputerName := "" //服务器名称
	sysOsName := ""       //操作系统
	sysOsArch := ""       //系统架构

	sysInfo, err := host.Info()

	if err == nil {
		sysComputerName = sysInfo.Hostname
		sysOsName = sysInfo.OS
		sysOsArch = sysInfo.KernelArch
	}

	goName := "GoLang"             //语言环境
	goVersion := runtime.Version() //版本
	gtime.Date()
	goStartTime := c.startTime //启动时间

	goRunTime := gtime.Now().Timestamp() - c.startTime.Timestamp() //运行时长（秒）
	goHome := runtime.GOROOT()                                     //安装路径
	goUserDir := ""                                                //项目路径

	curDir, err := os.Getwd()

	if err == nil {
		goUserDir = curDir
	}

	//服务器磁盘信息
	diskList := make([]disk.UsageStat, 0)
	diskInfo, err := disk.Partitions(true) //所有分区
	if err == nil {
		for _, p := range diskInfo {
			diskDetail, err := disk.Usage(p.Mountpoint)
			if err == nil {
				diskDetail.UsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", diskDetail.UsedPercent), 64)
				diskList = append(diskList, *diskDetail)
			}
		}
	}
	res = new(system.MonitorSearchRes)
	res = &system.MonitorSearchRes{
		"cpuNum":          cpuNum,
		"cpuUsed":         cpuUsed,
		"cpuAvg5":         gconv.String(cpuAvg5),
		"cpuAvg15":        gconv.String(cpuAvg15),
		"memTotal":        memTotal,
		"goTotal":         goTotal,
		"memUsed":         memUsed,
		"goUsed":          goUsed,
		"memFree":         memFree,
		"goFree":          goFree,
		"memUsage":        memUsage,
		"goUsage":         goUsage,
		"sysComputerName": sysComputerName,
		"sysOsName":       sysOsName,
		"sysComputerIp":   sysComputerIp,
		"sysOsArch":       sysOsArch,
		"goName":          goName,
		"goVersion":       goVersion,
		"goStartTime":     goStartTime,
		"goRunTime":       goRunTime,
		"goHome":          goHome,
		"goUserDir":       goUserDir,
		"diskList":        diskList,
	}
	return
}
