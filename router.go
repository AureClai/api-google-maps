package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rakyll/statik/fs"

	// prod only !!!
	_ "./statik"
)

var alreadyConnected bool = false

// Message the message
type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

func newRouter() *mux.Router {
	// Create the new router
	r := mux.NewRouter()

	// prod only !!
	// Statik files management
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	// Create the file server from statik
	h := http.FileServer(statikFS)

	r.HandleFunc("/ws", wsHandler)

	// Prod only !!
	r.PathPrefix("/").Handler(h).Methods("GET")

	return r
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	if alreadyConnected {
		fmt.Println("You are already connected on a page, close it and try again")
		http.Error(w, "", 403)
		return
	}

	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}

	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	defer conn.Close()

	alreadyConnected = true
	theModel.client = conn
	go messangeHandling()
	sendMessage(&Message{"settings change", theModel.Settings})
	sendAllPaths()
	write()
}
