package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	asciiart "github.com/romance-dev/ascii-art"
	gt "gopkg.gilang.dev/google-translate"
)

//go:embed ansi-shadow.flf
var font []byte

func init() {
	asciiart.RegisterFont("ansi-shadow", font)
}

func main() {
	figlet := asciiart.NewColorFigure("ID-EN", "ansi-shadow", "purple", true)
	figlet.Print()

	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	yellow := color.New(color.FgHiYellow, color.Bold).SprintFunc()

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print(yellow(">"))
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		// hapus newline
		input = strings.TrimSpace(input)

		translated, err := gt.ManualTranslate(input, "id", "en")
		if err != nil {
			panic(err)
		}

		fmt.Println(green("=>"), translated.Text)

	}

}
