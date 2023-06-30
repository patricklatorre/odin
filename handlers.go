package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/patricklatorre/odin/path"
)

// Downloads the Valheim server files into a new server dir within servers/
func create(name string) error {
	var (
		serverDir   = path.Relative("servers", name)
		steamcmdBin = path.Relative("steamcmd", "steamcmd.exe")
	)

	exists, err := path.Exists(serverDir)
	if err != nil {
		fmt.Printf("Could not check if %s exists\n", serverDir)
		return err
	}

	if exists {
		fmt.Printf("%s already exists\n", name)
		return nil
	}

	err = os.Mkdir(serverDir, 0755)
	if err != nil {
		fmt.Printf("Could not create the %s directory\n", serverDir)
		return err
	}

	fmt.Printf("Created directory: %s\n", name)

	// Setup command
	cmd := exec.Command(
		steamcmdBin,
		"+force_install_dir", serverDir,
		"+login", "anonymous",
		"+app_update", "896660",
		"validate", "+exit")

	// Pipe output to stdout
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	// Run command
	if err := cmd.Run(); err != nil {
		switch e := err.(type) {
		case *exec.Error:
			fmt.Println("Could not create the server")
			return err
		case *exec.ExitError:
			if e.ExitCode() != 7 { // Known bug
				fmt.Println("An error occurred while running steamcmd")
				return err
			}
		default:
			fmt.Println("An error occurred while running steamcmd")
			return err
		}
	}

	fmt.Printf("Created server \"%s\"\n", name)
	return nil
}

// Starts a server
func start(name string, port int, password string) error {
	var (
		serverDir = path.Relative("servers", name)
		serverBin = path.Relative("servers", name, "valheim_server.exe")
	)

	exists, err := path.Exists(serverDir)
	if err != nil {
		return err
	}

	if !exists {
		fmt.Printf("Server doesn't exist: %s\n", serverDir)
		os.Exit(1)
	}

	// Required by steamcmd
	os.Setenv("SteamAppId", "892970")

	cmd := exec.Command(
		serverBin,
		"-nographics",
		"-batchmode",
		"-name", name,
		"-world", name,
		"-port", strconv.Itoa(port),
		"-password", password,
		"-savedir", serverDir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Valheim process doesn't terminate automatically
	fmt.Println("Tip: Press CTRL+C to save and quit server")

	if err := cmd.Run(); err != nil {
		switch e := err.(type) {
		case *exec.Error:
			fmt.Println("Could not run the server")
			return err

		case *exec.ExitError:
			fmt.Println("An error occurred while running the server. Exit code:", e.ExitCode())
			return err

		default:
			fmt.Println("An error occurred while running the server")
			return err
		}
	}

	return nil
}

// Opens a server dir in the file explorer
func open(name string) error {
	serverDir := path.Relative("servers", name)

	exists, err := path.Exists(serverDir)
	if err != nil {
		fmt.Printf("Could not check if %s exists\n", serverDir)
		return err
	}

	if !exists {
		fmt.Println("Server doesn't exist:", serverDir)
		os.Exit(1)
	}

	// Windows-specific file explorer
	cmd := exec.Command("explorer", serverDir)
	_ = cmd.Run()

	return nil
}

// Prints the help screen
func help() {
	flag.Usage()
}

// Prints the version
func version() {
	fmt.Printf("Odin %s\n", Version)
}
