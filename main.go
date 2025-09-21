package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v3"
)

const LOCATION = "/sys/bus/platform/drivers/ideapad_acpi/VPC2004:00/conservation_mode"

func displayStatus() error {
	file, err := os.Open(LOCATION)
	if err != nil {
		return err
	}
	defer file.Close()


	var status int
	_, err = fmt.Fscanf(file, "%d", &status)
	if err != nil {
		return err
	}

	switch status {
	case 1:
		fmt.Println("battery protection enabled")
	case 0:
		fmt.Println("battery protection disabled")
	}
	return nil
}

func enableProtection() error {
	err := os.WriteFile(LOCATION, []byte("1"), 0644)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			return runAsRoot("-e")
		} else {
			return err
		}
	}

	return nil
}

func disableProtection() error {
	err := os.WriteFile(LOCATION, []byte("0"), 0644)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			return runAsRoot("-d")
		} else {
			return err
		}
	}
	return nil
}

func runAsRoot(arg string) error {
	fmt.Println("This operation requires root privileges, running again with sudo...")
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command("sudo", exe, arg)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func action(ctx context.Context, cmd *cli.Command) error {
	var err error

	if cmd.Bool("enable") {
		err = enableProtection()
	}

	if cmd.Bool("disable") {
		err = disableProtection()
	}

	if cmd.Bool("status") {
		err = displayStatus()
	}

	return err
}

func main() {
	cmd := &cli.Command{
		UseShortOptionHandling: true,
		EnableShellCompletion: true,
		Name:                   "battery-protection-toggle",
		Usage:                  "Control battery protection mode on Lenovo laptops.",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "enable", Aliases: []string{"e"}},
			&cli.BoolFlag{Name: "disable", Aliases: []string{"d"}},
			&cli.BoolFlag{Name: "status", Aliases: []string{"s"}, Value: true},
		},
		Action: action,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
