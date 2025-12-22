package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	url := "https://desustream.info/dstream/ondesu/hd/v3/index.php?id=ZWxDTkJkQ0VjVGJFZlpGd3FCZ2tsYm1TcDZDTVpicktobk51d2VvMDdNQT0"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	videoURL := findVideoSource(doc)
	if videoURL == "" {
		fmt.Println("Video source tidak ditemukan")
		return
	}

	fmt.Println("Video URL ditemukan:")
	fmt.Println(videoURL)

	downloadVideo(videoURL)
}

func findVideoSource(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "source" {
		for _, attr := range n.Attr {
			if attr.Key == "src" && strings.Contains(attr.Val, "googlevideo.com") {
				return attr.Val
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findVideoSource(c); result != "" {
			return result
		}
	}
	return ""
}

func downloadVideo(videoURL string) {
	resp, err := http.Get(videoURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	file, err := os.Create("video.mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.ReadFrom(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Download selesai: video.mp4")
}
