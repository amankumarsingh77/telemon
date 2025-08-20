package main

import (
	"encoding/json"
	"fmt"
	"github.com/amankumarsinghy77/telemon/constants"
	"net"
	"sync"
	"time"
)

type MetricStore interface {
	Store(metrics *constants.SystemMetrics)
	Query(hostname string, to, from time.Time) ([]*constants.SystemMetrics, error)
}

type CollectorServer struct {
	UDPAddr    string
	Storage    MetricStore
	metricChan chan *constants.SystemMetrics
	conn       *net.UDPConn
}

func NewCollectorServer(udpAddr string, storage MetricStore) *CollectorServer {
	return &CollectorServer{
		UDPAddr:    udpAddr,
		Storage:    storage,
		metricChan: make(chan *constants.SystemMetrics, 1000),
	}
}

func (s *CollectorServer) Start() error {
	addr, err := net.ResolveUDPAddr("udp", s.UDPAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve address %v", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s\n", err)
	}
	s.conn = conn
	fmt.Printf("started listening on %s\n", s.UDPAddr)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go s.metricProcessor(&wg)
	}
	s.listen()
	wg.Wait()
	return nil
}

func (s *CollectorServer) metricProcessor(wg *sync.WaitGroup) {
	defer wg.Done()
	for metrics := range s.metricChan {
		s.Storage.Store(metrics)
		s.checkAlerts(metrics)
	}
}

func (s *CollectorServer) checkAlerts(metrics *constants.SystemMetrics) {
	threshold := map[string]float64{
		"cpu_usage":  90.0,
		"mem_usage":  95.0,
		"disk_usage": 90.0,
		"avg_load":   5.0,
	}
	if metrics.CPUUsage > threshold["cpu_usage"] {
		s.triggerAlert(metrics.Hostname, "CPU", metrics.CPUUsage)
	}
	if metrics.MemUsage > threshold["mem_usage"] {
		s.triggerAlert(metrics.Hostname, "Memory", metrics.MemUsage)
	}
	if metrics.MemUsage > threshold["disk_usage"] {
		s.triggerAlert(metrics.Hostname, "Disk", metrics.DiskUsage)
	}
	if metrics.MemUsage > threshold["avg_load"] {
		s.triggerAlert(metrics.Hostname, "avg_load", metrics.AvgLoad)
	}
}

func (s *CollectorServer) triggerAlert(hostname, metric string, value float64) {
	// for now we will just print the alert
	fmt.Printf("ALERT: %s - %s usage is high: %.2f\n", hostname, metric, value)
}

func (s *CollectorServer) listen() {
	buffer := make([]byte, 1024)
	for {
		n, addr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading from UDP: %v\n", err)
			continue
		}
		var metrics constants.SystemMetrics
		if err := json.Unmarshal(buffer[:n], &metrics); err != nil {
			fmt.Printf("Error unmarshaling metrics from %s: %v\n", addr.String(), err)
			continue
		}

		s.metricChan <- &metrics
	}
}
