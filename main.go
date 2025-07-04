package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const LOCATION = "/sys/bus/platform/drivers/ideapad_acpi/VPC2004:00/conservation_mode"

func displayStatus() {
	file, err := os.Open(LOCATION)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var status int
	fmt.Fscanf(file, "%d", &status)

	switch status {
	case 1:
		fmt.Println("battery protection enabled")
	case 0:
		fmt.Println("battery protection disabled")
	}
}

func enableProtection() {
	err := os.WriteFile(LOCATION, []byte("1"), 0644)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			runAsRoot("-on")
		} else {
			log.Fatal(err)
		}
	}
}

func disableProtection() {
	err := os.WriteFile(LOCATION, []byte("0"), 0644)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			runAsRoot("-off")
		} else {
			log.Fatal(err)
		}
	}
}

func runAsRoot(arg string) {
	fmt.Println("This operation requires root privileges, running again with sudo...")
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("sudo", exe, arg)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	status := flag.Bool("status", false, "battery protection status")
	on := flag.Bool("on", false, "turn on battery protection")
	off := flag.Bool("off", false, "turn off battery protection")
	help := flag.Bool("help", false, "display help message")

	flag.Parse()

	switch {
	case *status == true:
		displayStatus()
	case *on == true:
		enableProtection()
	case *off == true:
		disableProtection()
	case *help == true:
		flag.Usage()
	default:
		flag.Usage()
		os.Exit(1)
	}
}
