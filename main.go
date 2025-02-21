package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

type SystemStats struct {
	Load           string `json:"cpu_load"`
	MemoryUsage    string `json:"memory_usage"`
	CPUTemperature string `json:"cpu_temperature"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	load, _ := cpu.Percent(0, false)
	memStats, _ := mem.VirtualMemory()
	cpuTemp, _ := getCPUTemperature()

	result := SystemStats{
		Load:           strconv.FormatFloat(load[0], 'f', 0, 64) + "%",
		MemoryUsage:    strconv.FormatFloat(memStats.UsedPercent, 'f', 0, 64) + "%",
		CPUTemperature: strconv.FormatFloat(cpuTemp, 'f', 0, 64) + "°C",
	}

	response := Response{
		Code:    200,
		Message: "Success",
		Payload: result,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getCPUTemperature() (float64, error) {
	tempFile := "/sys/class/thermal/thermal_zone0/temp"
	data, err := os.ReadFile(tempFile)
	if err != nil {
		return 0, err
	}

	tempStr := strings.TrimSpace(string(data))
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0, err
	}
	temp = temp / 1000.0
	return temp, nil
}

func myIPHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://ifconfig.me/ip")
	if err != nil {
		http.Error(w, "Unable to fetch IP", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	response := Response{
		Code:    200,
		Message: "Success",
		Payload: string(ip),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/my-ip", myIPHandler)
	http.ListenAndServe(":8000", nil)
}
