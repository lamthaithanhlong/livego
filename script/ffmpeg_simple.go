package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main_simple() {
	for {
		// Define the FFmpeg command
		cmd := exec.Command("ffmpeg", "-f", "avfoundation", "-framerate", "30", "-i", "0:4", "-af", "afftdn", "-s", "640x360", "-c:v", "libx264", "-preset", "veryfast", "-tune", "zerolatency", "-pix_fmt", "yuv420p", "-g", "30", "-b:v", "300k", "-bufsize", "300k", "-ar", "44100", "-b:a", "300k", "-c:a", "aac", "-profile:a", "aac_low", "-f", "flv", "rtmp://167.88.168.20:1935/live/test")

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
					fmt.Println("Error or freeze detected, handling...")
					// Handle the error or freeze condition
				}
			}
		}()

		// Wait for the FFmpeg command to exit
		err = cmd.Wait()
		if err != nil {
			fmt.Printf("FFmpeg exited with error: %v\n", err)
		} else {
			fmt.Println("FFmpeg exited without error.")
		}

		// Optional: handle the restart or termination here
	}
}
