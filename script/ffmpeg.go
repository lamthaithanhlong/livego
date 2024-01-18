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
	const (
		minInterval = 60  // minimum interval in seconds
		maxInterval = 360 // maximum interval in seconds
		step        = 30  // step to increase/decrease interval
	)
	restartInterval := 180 // starting interval

	for {
		// Define the FFmpeg command
		cmd := exec.Command("ffmpeg", "-f", "avfoundation", "-framerate", "30", "-i", "1:0", "-s", "640x360", "-c:v", "libx264", "-preset", "ultrafast", "-tune", "zerolatency", "-pix_fmt", "yuv420p", "-g", "30", "-b:v", "200k", "-c:a", "aac", "-f", "flv", "rtmp://167.88.168.20:1935/live/test")

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
		errorDetected := false
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Println(line)
				if strings.Contains(line, "error") || strings.Contains(line, "freeze") {
					fmt.Println("Error or freeze detected, restarting stream...")
					errorDetected = true
					cmd.Process.Kill()
				}
			}
		}()

		// Channel to receive exit status
		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		// Create a ticker with the adaptive interval
		ticker := time.NewTicker(time.Duration(restartInterval) * time.Second)
		defer ticker.Stop()

		// Wait for either the FFmpeg command to exit or the ticker to tick
		select {
		case <-ticker.C:
			fmt.Println("Restarting stream after interval...")
			cmd.Process.Kill()
		case err := <-done:
			if err != nil {
				fmt.Printf("FFmpeg exited with error: %v\n", err)
			} else {
				fmt.Println("FFmpeg exited without error.")
			}
		}

		// Adjust restartInterval based on error detection
		if errorDetected {
			restartInterval = max(minInterval, restartInterval-step)
		} else {
			restartInterval = min(maxInterval, restartInterval+step)
		}
	}
}

// Helper functions for min and max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
