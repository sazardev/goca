package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func gocaBinary() string {
	if runtime.GOOS == "windows" {
		return "goca.exe"
	}
	return "goca"
}

func main() {
	goca := gocaBinary()

	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "goca-test-*")
	if err != nil {
		fmt.Println("Error creating temp dir:", err)
		os.Exit(1)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			fmt.Printf("Error cleaning up temp directory: %v\n", err)
		}
	}()

	// Initialize project
	cmd := exec.Command(goca, "init", "testproject", "--module", "github.com/test/testproject")
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error initializing project:", err)
		os.Exit(1)
	}

	// Create entity
	cmd = exec.Command(goca, "entity", "User", "--fields", "name:string,email:string,age:int")
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error creating entity:", err)
		os.Exit(1)
	}

	// Read and display the file content
	userFilePath := filepath.Join(tmpDir, "internal", "domain", "user.go")
	fmt.Println("Contents of", userFilePath, ":")
	content, err := os.ReadFile(userFilePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	fmt.Println(string(content))
}
