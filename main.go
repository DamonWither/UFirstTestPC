package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// Requirements defines the minimum hardware requirements.
type Requirements struct {
	Cores         int     `json:"cores"`
	ClockSpeedGHz float64 `json:"clock_speed_ghz"`
	MemoryGB      int     `json:"memory_gb"`
	DiskSpaceGB   int     `json:"disk_space_gb"`
}

// PCConfig defines the configuration of a PC.
type PCConfig struct {
	Cores         int     `json:"cores"`
	ClockSpeedGHz float64 `json:"clock_speed_ghz"`
	MemoryGB      int     `json:"memory_gb"`
	DiskSpaceGB   int     `json:"disk_space_gb"`
}

// Result defines the structure of the response.
type Result struct {
	MeetsRequirements bool   `json:"meets_requirements"`
	Message           string `json:"message"`
	MessageBecause    string `json:"message_because"`
	Message2          string `json:"message2"`
	Message3          string `json:"message3"`
	Message4          string `json:"message4"`
	Message5          string `json:"message5"`
}

var tmpl = template.Must(template.New("form").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>UFirst - Проверка конфигурации ПК</title>
</head>
<body>
	<h1>UFirst - Проверка конфигурации ПК</h1>
	<h2>{{.Message}}</h2>
	<h2>{{.MessageBecause}}</h2>
	<h2>{{.Message2}}</h2>
	<h2>{{.Message3}}</h2>
	<h2>{{.Message4}}</h2>
	<h2>{{.Message5}}</h2>
</body>
</html>
`))

func getSystemConfig() PCConfig {
	cores, err := cpu.Counts(true)
	if err != nil {
		log.Println("Error getting CPU cores:", err)
		cores = 0
	}

	clockSpeed := 0.0
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		clockSpeed = float64(cpuInfo[0].Mhz) / 1000.0
	} else if err != nil {
		log.Println("Error getting CPU clock speed:", err)
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting memory info:", err)
	}
	memory := int(vm.Total / (1024 * 1024 * 1024))

	diskInfo, err := disk.Usage("/")
	if err != nil {
		log.Println("Error getting disk space:", err)
	}
	diskSpace := int(diskInfo.Free / (1024 * 1024 * 1024))

	log.Printf("Detected system configuration: Cores=%d, ClockSpeed=%.2fGHz, Memory=%dGB, DiskSpace=%dGB\n",
		cores, clockSpeed, memory, diskSpace)

	return PCConfig{
		Cores:         cores,
		ClockSpeedGHz: clockSpeed,
		MemoryGB:      memory,
		DiskSpaceGB:   diskSpace,
	}
}

func checkSystemHandler(w http.ResponseWriter, r *http.Request) {
	required := Requirements{
		Cores:         4,
		ClockSpeedGHz: 2.1,
		MemoryGB:      6,
		DiskSpaceGB:   20,
	}

	config := getSystemConfig()

	result := Result{}
	if config.Cores >= required.Cores &&
		config.ClockSpeedGHz >= required.ClockSpeedGHz &&
		config.MemoryGB >= required.MemoryGB &&
		config.DiskSpaceGB >= required.DiskSpaceGB {
		result.MeetsRequirements = true
		result.Message = "Ваш ПК удовлетворяет требованиям."
	} else {
		result.MessageBecause = "Это произошло из-за:"
		if config.Cores < required.Cores {
			result.Message2 = "Недостаточно ядер процессора"
		}
		if config.ClockSpeedGHz < required.ClockSpeedGHz {
			result.Message3 = "Недостаточно тактовой частоты"
		}
		if config.MemoryGB < required.MemoryGB {
			result.Message4 = "Недостаточно оперативной памяти"
		}
		if config.DiskSpaceGB < required.DiskSpaceGB {
			result.Message5 = "Недостаточно места диске"
		}
		result.MeetsRequirements = false
		result.Message = "Ваш ПК не удовлетворяет требованиям."
	}

	tmpl.Execute(w, result)
}

func main() {
	http.HandleFunc("/", checkSystemHandler)

	port := ":8080"
	log.Println("Server is running on http://localhost" + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
