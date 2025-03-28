package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var typeName = "WSBenchLogs"

func statsAnalysis() {
	data, err := parseStatsLog("./benchmark/stats.txt")
	if err != nil {
		fmt.Println("Error reading stats log:", err)
		return
	}

	plotCPUStats(data)
	plotMemoryStats(data)
}

type StatsData struct {
	Time   float64
	CPU    float64
	Memory float64
}

func parseStatsLog(filename string) ([]StatsData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := csv.NewReader(file)
	scanner.FieldsPerRecord = -1

	var data []StatsData
	reCPU := regexp.MustCompile(`([0-9.]+)%`)
	timeCounter := 0.0

	lines, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	entries := strings.Split(string(lines), "\n")
	for i := 1; i < len(entries); i += 2 {
		fields := strings.Fields(entries[i])

		cpuMatch := reCPU.FindStringSubmatch(fields[2])
		if len(cpuMatch) < 2 {
			continue
		}
		cpu, _ := strconv.ParseFloat(cpuMatch[1], 64)

		memParts := strings.Fields(fields[3])
		mem, _ := parseMemory(memParts[0])

		data = append(data, StatsData{Time: timeCounter, CPU: cpu, Memory: mem})
		timeCounter += 1
	}

	return data, nil
}

func parseMemory(memStr string) (float64, error) {
	var unitMap = map[string]float64{
		"KiB": 1, "MiB": 1024, "GiB": 1024 * 1024,
	}
	for unit, factor := range unitMap {
		if strings.HasSuffix(memStr, unit) {
			value, err := strconv.ParseFloat(strings.TrimSuffix(memStr, unit), 64)
			if err != nil {
				return 0, err
			}
			return value * factor, nil
		}
	}
	return 0, fmt.Errorf("unknown unit: %s", memStr)
}

func plotCPUStats(data []StatsData) {
	p := plot.New()
	p.Title.Text = "CPU Usage Over Time"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "CPU (%)"

	ptsCPU := make(plotter.XYs, len(data))

	for i, d := range data {
		ptsCPU[i].X = d.Time
		ptsCPU[i].Y = d.CPU
	}

	lineCPU, _ := plotter.NewLine(ptsCPU)
	lineCPU.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	p.Add(lineCPU)
	p.Legend.Add("CPU (%)", lineCPU)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, "cpu_usage.png"); err != nil {
		fmt.Println("Error saving CPU plot:", err)
	}
}

func plotMemoryStats(data []StatsData) {
	p := plot.New()
	p.Title.Text = "Memory Usage Over Time"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Memory (MiB)"

	ptsMem := make(plotter.XYs, len(data))

	for i, d := range data {
		ptsMem[i].X = d.Time
		ptsMem[i].Y = d.Memory / 1024 // MiB 表示
	}

	lineMem, _ := plotter.NewLine(ptsMem)
	lineMem.Color = color.RGBA{B: 255, A: 255}
	p.Add(lineMem)
	p.Legend.Add("Memory (MiB)", lineMem)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, "memory_usage.png"); err != nil {
		fmt.Println("Error saving Memory plot:", err)
	}
}

type LogEntry struct {
	Event     string
	UserID    string
	Timestamp time.Time
}

func parseLogLine(line string) (LogEntry, error) {
	parts := strings.Split(line, ", ")
	if len(parts) != 3 {
		return LogEntry{}, fmt.Errorf("invalid log format: %s", line)
	}
	timeParsed, err := time.Parse("15:04:05.000", parts[2])
	if err != nil {
		return LogEntry{}, err
	}
	return LogEntry{Event: parts[0], UserID: parts[1], Timestamp: timeParsed}, nil
}

func readLogFile(filename string) (map[string][]LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	logs := make(map[string][]LogEntry)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry, err := parseLogLine(scanner.Text())
		if err != nil {
			return nil, err
		}
		logs[entry.UserID] = append(logs[entry.UserID], entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func calculateLatency(filename string, clientLogs, serverLogs map[string][]LogEntry) {
	var clientToServerLatencies []time.Duration
	var serverToClientLatencies []time.Duration

	for userID, clientEntries := range clientLogs {
		serverEntries, exists := serverLogs[userID]
		if !exists {
			continue
		}

		var requestTimes []time.Time
		var responseTimes []time.Time

		for _, entry := range clientEntries {
			if entry.Event == "send" {
				requestTimes = append(requestTimes, entry.Timestamp)
			} else if entry.Event == "comes" {
				responseTimes = append(responseTimes, entry.Timestamp)
			}
		}

		for _, entry := range serverEntries {
			if entry.Event == "comes" && len(requestTimes) > 0 {
				clientToServerLatencies = append(clientToServerLatencies, entry.Timestamp.Sub(requestTimes[0]))
				requestTimes = requestTimes[1:]
			} else if entry.Event == "send" && len(responseTimes) > 0 {
				if responseTimes[0].Sub(entry.Timestamp) < 0 {
					fmt.Println(entry)
				}
				serverToClientLatencies = append(serverToClientLatencies, responseTimes[0].Sub(entry.Timestamp))
				responseTimes = responseTimes[1:]
			}
		}
	}

	sort.Slice(clientToServerLatencies, func(i, j int) bool { return clientToServerLatencies[i] < clientToServerLatencies[j] })
	sort.Slice(serverToClientLatencies, func(i, j int) bool { return serverToClientLatencies[i] < serverToClientLatencies[j] })

	// fmt.Println("Client to Server Latency results:")
	// for i, latency := range clientToServerLatencies {
	// 	fmt.Printf("%d: %s\n", i+1, latency)
	// }

	// fmt.Println("Server to Client Latency results:")
	// for i, latency := range serverToClientLatencies {
	// 	fmt.Printf("%d: %s\n", i+1, latency)
	// }

	writeLatencyToFile(fmt.Sprintf("client_to_server_latency_%s.txt", filename), clientToServerLatencies)
	writeLatencyToFile(fmt.Sprintf("server_to_client_latency_%s.txt", filename), serverToClientLatencies)
}

func writeLatencyToFile(filename string, latencies []time.Duration) {
	if len(latencies) == 0 {
		return
	}

	minLatency := latencies[0]
	maxLatency := latencies[len(latencies)-1]
	mean := calcMean(latencies)
	stdDev := calcStdDev(latencies, mean)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "Min Latency: %s\n", minLatency)
	fmt.Fprintf(file, "Max Latency: %s\n", maxLatency)
	fmt.Fprintf(file, "Mean: %s\n", mean)
	fmt.Fprintf(file, "Standard Deviation: %s\n", stdDev)
	fmt.Fprintln(file, "Latency Data:")

	for _, latency := range latencies {
		fmt.Fprintln(file, latency)
	}
}

func calcMean(latencies []time.Duration) time.Duration {
	sum := time.Duration(0)
	for _, latency := range latencies {
		sum += latency
	}
	return sum / time.Duration(len(latencies))
}

func calcStdDev(latencies []time.Duration, mean time.Duration) time.Duration {
	var sumSquares float64
	for _, latency := range latencies {
		diff := float64(latency - mean)
		sumSquares += diff * diff
	}
	variance := sumSquares / float64(len(latencies))
	return time.Duration(math.Sqrt(variance))
}

func main() {
	fmt.Println(typeName)
	clientLogs, err := readLogFile(fmt.Sprintf("./benchmark/%s.txt", typeName))
	if err != nil {
		fmt.Println("Error reading client log:", err)
		return
	}

	serverLogs, err := readLogFile(fmt.Sprintf("./app/%s.txt", typeName))
	if err != nil {
		fmt.Println("Error reading server log:", err)
		return
	}

	calculateLatency("post", clientLogs, serverLogs)

	clientLogs, err = readLogFile(fmt.Sprintf("./benchmark/%sImage.txt", typeName))
	if err != nil {
		fmt.Println("Error reading client log:", err)
		return
	}

	serverLogs, err = readLogFile(fmt.Sprintf("./app/%sImage.txt", typeName))
	if err != nil {
		fmt.Println("Error reading server log:", err)
		return
	}

	calculateLatency("image", clientLogs, serverLogs)

	statsAnalysis()
}
