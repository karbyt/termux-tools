package main

import (
	_ "embed"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/fatih/color"
	asciiart "github.com/romance-dev/ascii-art"
)

//go:embed ansi-shadow.flf
var font []byte

type Device struct {
	Name string
	IP   string
}

func init() {
	asciiart.RegisterFont("ansi-shadow", font)
}

// fungsi ping menggunakan command system
func ping(ip string) bool {
	// -c 1 : ping 1 kali
	// -W 1 : timeout 1 detik (Linux/Termux)
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)
	err := cmd.Run()
	return err == nil
}

func waitExit() {
	fmt.Println("\nTekan ENTER untuk keluar...")
	fmt.Scanln()
}

func main() {

	devices := []Device{
		{"Ringlight  ", "192.168.1.111"},
		{"Lampu Mejo ", "192.168.1.112"},
		{"Jam        ", "192.168.1.103"},
		{"Kipas Cilik", "192.168.1.114"},
		{"Server     ", "192.168.1.115"},
		{"Relay      ", "192.168.1.116"},
		{"OPPO A3S   ", "192.168.1.4"},
		{"Server     ", "192.168.1.100"},
	}

	results := make([]bool, len(devices))
	var wg sync.WaitGroup

	figlet := asciiart.NewColorFigure("PING", "ansi-shadow", "purple", true)
	figlet.Print()

	fmt.Println("  Mengecek status perangkat...")
	fmt.Println("  ------------------------------")

	start := time.Now()

	for i, dev := range devices {
		wg.Add(1)

		go func(index int, d Device) {
			defer wg.Done()
			results[index] = ping(d.IP)
		}(i, dev)
	}

	wg.Wait()

	elapsed := time.Since(start)

	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	red := color.New(color.FgRed, color.Bold).SprintFunc()

	for i, dev := range devices {
		statusText := ""
		if results[i] {
			statusText = green("● UP")
		} else {
			statusText = red("✖ DOWN")
		}

		fmt.Printf("  %-12s (%-13s) : %s\n", dev.Name, dev.IP, statusText)
	}

	fmt.Println("  ------------------------------")
	fmt.Printf("  Selesai dalam %v\n", elapsed)

	waitExit()
}
