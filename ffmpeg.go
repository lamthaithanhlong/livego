package main

import (
	"log"
	"os/exec"
	"time"
)

func main() {
	for {
		// Define the FFmpeg command
		cmd := exec.Command("ffmpeg", "-f", "avfoundation", "-framerate", "30", "-i", "1:0", "-s", "854x480", "-c:v", "libx264", "-preset", "ultrafast", "-tune", "zerolatency", "-pix_fmt", "yuv420p", "-g", "5", "-b:v", "200k", "-c:a", "aac", "-f", "flv", "rtmp://127.0.0.1:1935/live/test")

		// Start the FFmpeg command
		err := cmd.Start()
		if err != nil {
			log.Fatalf("Failed to start FFmpeg: %v", err)
		}

		// Wait for 5 seconds
		time.Sleep(10 * time.Minute)

		// Stop the FFmpeg command
		err = cmd.Process.Kill()
		if err != nil {
			log.Fatalf("Failed to kill FFmpeg: %v", err)
		}

		// Wait a little before restarting
		time.Sleep(1 * time.Second)
	}
}
