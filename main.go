package main

import (
	"fmt"
	"log"
	"os"
)

const LOCATION = "/sys/bus/platform/drivers/ideapad_acpi/VPC2004:00/conservation_mode"

func usage() {
	message := `usage: battery-protection-toggle [operation]
    operations:
        -h  Display this message.
        -s  Display battery protection status.
        -e  Enable battery protection.
        -d  Disable battery protection.
    `
	fmt.Printf(message)
}

func status() {
	file, err := os.Open(LOCATION)
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

func enable_protection() {
	err := os.WriteFile(LOCATION, []byte("1"), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func disable_protection() {
	err := os.WriteFile(LOCATION, []byte("0"), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "-s":
		status()
	case "-e":
		enable_protection()
	case "-d":
		disable_protection()
	case "-h":
		usage()
	default:
		log.Print("unknown option")
	}
}
