package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	for {
		// Define the FFmpeg command
		cmd := exec.Command("ffmpeg", "-f", "avfoundation", "-framerate", "30", "-i", "1:0", "-s", "640x360", "-c:v", "libx264", "-preset", "ultrafast", "-tune", "zerolatency", "-pix_fmt", "yuv420p", "-g", "30", "-b:v", "200k", "-c:a", "aac", "-f", "flv", "rtmp://167.88.168.20:1935/live/test")
		// cmd := exec.Command("ffmpeg", "-i", "rtmp://localhost/live/livestream", "-c", "copy", "-f", "flv", "rtmp://167.88.168.20:1935/live/test")

		// Create stderr pipe
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatalf("Failed to create stderr pipe: %v", err)
		}

		// Start FFmpeg
		err = cmd.Start()
		fmt.Println("Stream started")
		if err != nil {
			log.Fatalf("Failed to start FFmpeg: %v", err)
		}

		// Scanner for stderr
		scanner := bufio.NewScanner(stderr)
		go func() {
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Println(line)
				if strings.Contains(line, "error") || strings.Contains(line, "freeze") {
					fmt.Println("Error or freeze detected, restarting stream...")
					cmd.Process.Kill()
				}
			}
		}()

		// Channel to receive exit status
		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		// Create a ticker for 180 seconds
		ticker := time.NewTicker(180 * time.Second)
		defer ticker.Stop()

		// Wait for either the FFmpeg command to exit or the ticker to tick
		select {
		case <-ticker.C:
			fmt.Println("180 seconds passed, restarting stream...")
			cmd.Process.Kill()
		case err := <-done:
			if err != nil {
				fmt.Printf("FFmpeg exited with error: %v\n", err)
			} else {
				fmt.Println("FFmpeg exited without error.")
			}
		}
	}
}
