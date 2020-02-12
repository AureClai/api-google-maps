package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
)

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func main() {
	router := newRouter()
	theModel = InitializeModel()
	fmt.Println("Model initialized")

	// prod only
	go open("http://localhost:8080")
	http.ListenAndServe(":8080", router)
	return
}
