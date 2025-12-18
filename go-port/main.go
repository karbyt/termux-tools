package main

import (
	_ "embed"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
	asciiart "github.com/romance-dev/ascii-art"
	"github.com/schollz/progressbar/v3"
)

//go:embed ansi-shadow.flf
var font []byte

func init() {
	asciiart.RegisterFont("ansi-shadow", font)
}

const (
	subnet      = "192.168.1."
	timeout     = 120 * time.Millisecond
	workerCount = 400
)

var ipRanges = [][2]int{
	{1, 20},    // DHCP
	{100, 120}, // Static / IoT
}

var popularPorts = []int{
	21, 22, 23, 25, 53,
	80, 110, 139, 143,
	443, 445,

	3000, 3001, 3002,

	4000, 5000,
	5432, 5900, 5173,

	8000, 8008, 8080, 8081,
}

type Job struct {
	Host string
	Port int
}

type Result struct {
	Host string
	Port int
}

// ---------- NETWORK CORE ----------

func scanPort(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial("tcp4", address)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func hostAlive(host string) bool {
	checkPorts := []int{22, 80, 443}
	for _, p := range checkPorts {
		if scanPort(host, p) {
			return true
		}
	}
	return false
}

// ---------- WORKER ----------

func worker(
	jobs <-chan Job,
	results chan<- Result,
	wg *sync.WaitGroup,
	bar *progressbar.ProgressBar,
) {
	defer wg.Done()

	for job := range jobs {
		if scanPort(job.Host, job.Port) {
			results <- Result(job)
		}
		bar.Add(1)
	}
}

func waitForExit() {
	color.Yellow("\nPress ENTER to exit...")
	fmt.Scanln()
}

// ---------- MAIN ----------

func main() {
	figlet := asciiart.NewColorFigure("PORT", "ansi-shadow", "purple", true)
	figlet.Print()
	color.Cyan("Subnet: 192.168.1.0/24")
	color.Cyan("Ranges: 1–20, 100–120\n")

	// calculate job count
	hostCount := 0
	for _, r := range ipRanges {
		hostCount += r[1] - r[0] + 1
	}
	totalJobs := hostCount * len(popularPorts)

	bar := progressbar.NewOptions(
		totalJobs,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetDescription("Scanning"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(20),
		progressbar.OptionThrottle(120*time.Millisecond),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetRenderBlankState(true),
	)

	jobs := make(chan Job, 1000)
	results := make(chan Result, 200)

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg, bar)
	}

	// producer
	go func() {
		for _, r := range ipRanges {
			for i := r[0]; i <= r[1]; i++ {
				host := fmt.Sprintf("%s%d", subnet, i)

				if !hostAlive(host) {
					bar.Add(len(popularPorts))
					continue
				}

				for _, port := range popularPorts {
					jobs <- Job{Host: host, Port: port}
				}
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	openPorts := make(map[string][]int)

	for res := range results {
		openPorts[res.Host] = append(openPorts[res.Host], res.Port)
	}

	// ---------- OUTPUT ----------
	color.Green("\n✅ Scan finished\n")

	if len(openPorts) == 0 {
		color.Red("No open ports found.")
		return
	}

	color.Yellow("📋 Open Ports:\n")

	for host, ports := range openPorts {
		color.Cyan("Host: %s", host)
		for _, p := range ports {
			color.Green("  └─ Port %d OPEN", p)
		}
		fmt.Println()
	}
	waitForExit()
}
