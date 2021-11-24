package main

import (
	"fmt"
	"github.com/ysksuzuki/containers-from-scratch/pkg/timens"
	"os"
	"os/exec"
	"runtime"
)

// go run main.go run <cmd> <args>
func main() {
	switch os.Args[1] {
	case "run":
		must(run())
	default:
		panic("help")
	}
}

func run() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	currentNs, err := timens.Get()
	if err != nil {
		return fmt.Errorf("failed to get the current time NS: %w", err)
	}
	defer timens.Set(currentNs)

	nsHandle, err := timens.New()
	if err != nil {
		return fmt.Errorf("failed to create time namespace: %w", err)
	}
	defer nsHandle.Close()

	cmd := exec.Command("bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
