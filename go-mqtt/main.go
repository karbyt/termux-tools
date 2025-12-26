package main

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	asciiart "github.com/romance-dev/ascii-art"
)

//go:embed ansi-shadow.flf
var font []byte

func init() {
	asciiart.RegisterFont("ansi-shadow", font)
}

func waitExit() {
	fmt.Println("\nTekan ENTER untuk keluar...")
	fmt.Scanln()
}

func main() {
	fig := asciiart.NewFigure("MQTT", "ansi-shadow", true)
	fig.Print()

	// Konfigurasi broker HiveMQ
	broker := "mqtts://15a5b434042f43b6a53b245e205f68cd.s1.eu.hivemq.cloud:8883"
	clientID := "go-mqtt-client"
	topic := "cmnd/kipaspaiyan/POWER"
	username := "Sandemo"
	password := "Sandemo787898"

	// Opsi client
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)
	opts.SetUsername(username)
	opts.SetPassword(password)

	// Optional: callback jika koneksi berhasil
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("Connected to HiveMQ broker")
	}

	// Optional: callback jika koneksi terputus
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Println("Connection lost:", err)
	}

	// Buat client
	client := mqtt.NewClient(opts)

	// Connect ke broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// Publish pesan
	token := client.Publish(topic, 1, true, "on")
	token.Wait()
	fmt.Println("Published message: online")

	time.Sleep(1000 * time.Millisecond)

	token = client.Publish(topic, 1, true, "off")
	token.Wait()
	fmt.Println("Published message: offline")

	// Tunggu sebentar sebelum disconnect
	time.Sleep(1 * time.Second)

	client.Disconnect(250)
	fmt.Println("Disconnected")

	waitExit()
}
