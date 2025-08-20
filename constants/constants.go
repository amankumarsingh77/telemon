package constants

type SystemMetrics struct {
	Timestamp   int64   `json:"timestamp"`
	Hostname    string  `json:"hostname"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemUsage    float64 `json:"mem_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	AvgLoad     float64 `json:"avg_load"`
	NetByteSent uint64  `json:"net_byte_sent"`
	NetByteRecv uint64  `json:"net_byte_recv"`
}
