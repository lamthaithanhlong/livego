package main

import (
	"fmt"
	"os/exec"
)

func main_audio() {
	// Get the default output device
	outputDeviceCmd := exec.Command("SwitchAudioSource", "-c", "-t", "output")
	outputDevice, err := outputDeviceCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Default Output Device: %s\n", outputDevice)

	// Get the default input device
	inputDeviceCmd := exec.Command("SwitchAudioSource", "-c", "-t", "input")
	inputDevice, err := inputDeviceCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Default Input Device: %s\n", inputDevice)
}
