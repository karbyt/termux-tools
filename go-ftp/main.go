package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"

	asciiart "github.com/romance-dev/ascii-art"
)

//go:embed ansi-shadow.flf
var font []byte

func init() {
	asciiart.RegisterFont("ansi-shadow", font)
}

func main() {
	figlet := asciiart.NewColorFigure("FTP", "ansi-shadow", "purple", true)
	figlet.Print()
	fmt.Println("FTP server starting on :2121")
	fmt.Println("Root: /sdcard")
	fmt.Println("Tekan ENTER untuk berhenti")

	cmd := exec.Command(
		"busybox",
		"tcpsvd",
		"-vE",
		"0.0.0.0",
		"2121",
		"busybox",
		"ftpd",
		"-w",
		"/sdcard",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
