package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

func main() {
	for {
		// Define the FFmpeg command
		cmd := exec.Command("ffmpeg", "-f", "avfoundation", "-framerate", "30", "-i", "1:1", "-s", "854x480", "-c:v", "libx264", "-preset", "ultrafast", "-tune", "zerolatency", "-pix_fmt", "yuv420p", "-g", "5", "-b:v", "200k", "-c:a", "aac", "-f", "flv", "rtmp://167.88.168.20:1935/live/test")

		// Create a pipe to read from standard error (stderr)
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatalf("Failed to create stderr pipe: %v", err)
		}

		// Start the FFmpeg command
		err = cmd.Start()
		fmt.Println("Stream started")
		if err != nil {
			log.Fatalf("Failed to start FFmpeg: %v", err)
		}

		// Use a scanner to read error output line by line
		scanner := bufio.NewScanner(stderr)
		go func() {
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Println(line) // Print FFmpeg's stderr to our stdout
				if strings.Contains(line, "error") || strings.Contains(line, "freeze") {
					fmt.Println("Error or freeze detected, restarting stream...")
					cmd.Process.Kill()
				}
			}
		}()

		// Wait for the FFmpeg command to exit or be killed
		err = cmd.Wait()
		if err != nil {
			fmt.Printf("FFmpeg exited with error: %v\n", err)
		} else {
			fmt.Println("FFmpeg exited without error.")
		}

		// Wait a little before restarting
		time.Sleep(1 * time.Second)
	}
}
