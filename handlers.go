package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/patricklatorre/odin/path"
)

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

func start(name string, port int, password string) error {
	fmt.Println("Start:", name, port, password)
	return nil
}

func open(name string) error {
	fmt.Println("Open:", name)
	return nil
}

func help() {
	flag.Usage()
	os.Exit(0)
}
