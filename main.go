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
    help  Display this message.
    status  Display battery protection status.
    on  Enable battery protection.
    off  Disable battery protection.
`
	fmt.Printf(message)
}

func status() {
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
	case "status":
		status()
	case "on":
		enable_protection()
	case "off":
		disable_protection()
	case "help":
		usage()
	default:
        fmt.Fprintln(os.Stderr, "unknown option")
	}
}
