package main

import (
	"net/http"
	"os/exec"
	"runtime"

	_ "./statik"
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

	go open("http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
