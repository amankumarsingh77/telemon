package main

import (
	"encoding/json"
	"fmt"
	"github.com/amankumarsinghy77/telemon/constants"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	util_net "github.com/shirou/gopsutil/v4/net"
	"net"
	"os"
	"time"
)

type TelemetryAgent struct {
	Conn         net.Conn
	ServerAddr   string
	Hostname     string
	PollInterval time.Duration
}

func NewTelemetry(serverAdd string, pollInterval time.Duration) *TelemetryAgent {
	hostname, _ := os.Hostname()
	return &TelemetryAgent{
		ServerAddr:   serverAdd,
		PollInterval: pollInterval,
		Hostname:     hostname,
	}
}

func (t *TelemetryAgent) Connect() error {
	conn, err := net.Dial("udp", t.ServerAddr)
	if err != nil {
		return fmt.Errorf("failed to connect : %v", err)
	}
	t.Conn = conn
	return nil
}

func (t *TelemetryAgent) CollectMetrics() (*constants.SystemMetrics, error) {
	cpuUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get cpu usage %v", err)
	}
	memUsage, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory usage %v", err)
	}
	diskUsage, err := disk.Usage("/")
	if err != nil {
		return nil, fmt.Errorf("failed to get disk usage %v", err)
	}
	avgLoad, err := load.Avg()
	if err != nil {
		return nil, fmt.Errorf("failed to get avg load %v", err)
	}
	netBytes, err := util_net.IOCounters(false)
	if err != nil {
		return nil, fmt.Errorf("failed to network I/O %v", err)
	}
	return &constants.SystemMetrics{
		Timestamp:   time.Now().UnixNano(),
		Hostname:    t.Hostname,
		CPUUsage:    cpuUsage[0],
		MemUsage:    memUsage.UsedPercent,
		DiskUsage:   diskUsage.UsedPercent,
		AvgLoad:     avgLoad.Load1,
		NetByteSent: netBytes[0].BytesSent,
		NetByteRecv: netBytes[0].BytesRecv,
	}, nil
}

func (t *TelemetryAgent) SendMetrics(metrics *constants.SystemMetrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics data : %v", err)
	}
	_, err = t.Conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send metrics : %v", err)
	}
	return nil
}

func (t *TelemetryAgent) Run() {
	ticker := time.NewTicker(t.PollInterval)
	defer ticker.Stop()
	for range ticker.C {
		metrics, err := t.CollectMetrics()
		if err != nil {
			fmt.Printf("failed to get metrics : %v", err)
		}
		if err = t.SendMetrics(metrics); err != nil {
			fmt.Printf("failed to send metrics : %v", err)
		}
	}
}
