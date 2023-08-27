//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "cli-pasta", ".")
	return cmd.Run()
}

func BuildDeamon() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "pasta-deamon", "./pasta-deamon")
	return cmd.Run()
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	return os.Rename("./cli-pasta", "/usr/bin/cli-pasta")
}

func InstallDeamon() error {
	mg.Deps(BuildDeamon)
	fmt.Println("Installing...")
	os.Rename("./pasta-deamon/pasta-deamon", "/usr/bin/pasta-deamon")
	return os.Rename("./pasta-deamon/pasta-deamon", "/usr/bin/pasta-deamon")
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "get", "github.com/stretchr/piglatin")
	return cmd.Run()
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("cli-pasta")
    os.RemoveAll("./pasta-deamon/pasta-deamon")
}

func Test() {
}
