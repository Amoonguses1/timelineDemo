package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

var typeName = "WSBenchLogs"

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

func calculateLatency(clientLogs, serverLogs map[string][]LogEntry) {
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

	writeLatencyToFile("client_to_server_latency.txt", clientToServerLatencies)
	writeLatencyToFile("server_to_client_latency.txt", serverToClientLatencies)
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

	calculateLatency(clientLogs, serverLogs)
}
