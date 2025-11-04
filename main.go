package main

import (
	"fmt"
	"os"

	"github.com/hybridgroup/mjpeg"
)

func main() {
	if len(os.Args) < 6 {
		fmt.Println("How to run:\n\tvideo-description [camera ID] [host:port] [model path] [projector path] [prompt text]")
		return
	}

	// parse args
	deviceID := os.Args[1]
	host := os.Args[2]
	modelPath := os.Args[3]
	projectorPath := os.Args[4]
	promptText := os.Args[5]

	stream := mjpeg.NewStream()

	go startVideoCapture(deviceID, stream)
	go startVLM(modelPath, projectorPath, promptText)

	fmt.Println("Capturing. Point your browser to", host)

	startWebServer(host, stream)
}
