package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kbinani/screenshot"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // allow browser
		},
	}
)

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.Handle("/", http.FileServer(http.Dir("./public")))

	go screenshotLoop()

	log.Println("Server jalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	log.Println("Client terhubung")

	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
		log.Println("Client terputus")
	}()

	// keep connection alive
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

func screenshotLoop() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		clientsMu.Lock()
		if len(clients) == 0 {
			clientsMu.Unlock()
			continue
		}
		clientsMu.Unlock()

		imgBytes, err := captureJPEG()
		if err != nil {
			log.Println("Screenshot error:", err)
			continue
		}

		broadcast(imgBytes)
	}
}

func captureJPEG() ([]byte, error) {
	scale := 1.5 // samakan dengan hyprctl monitors
	b := screenshot.GetDisplayBounds(0)

	bounds := image.Rect(
		int(float64(b.Min.X)*scale),
		int(float64(b.Min.Y)*scale),
		int(float64(b.Max.X)*scale),
		int(float64(b.Max.Y)*scale),
	)
	// bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{
		Quality: 60, // penting untuk realtime
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func broadcast(data []byte) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn := range clients {
		err := conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			conn.Close()
			delete(clients, conn)
		}
	}
}
