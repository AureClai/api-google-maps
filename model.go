package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

// Model struct
type Model struct {
	client   *websocket.Conn
	send     chan Message
	Settings *Settings
	mu       sync.Mutex
}

// Parameter struct
type Settings struct {
	APIKey string `json:"api-key"`
	Proxy  string `json:"proxy"`
}

// new empty settings
func getNewEmptySettings() *Settings {
	return &Settings{
		APIKey: "",
		Proxy:  "",
	}
}

var theModel Model

func InitializeModel() Model {
	jsonFile := GetSettingsFile()
	settings := getNewEmptySettings()
	if jsonFile != nil {
		settings = readSettingsFromFile(jsonFile)
	}
	model := Model{
		Settings: settings,
		client:   nil,
		send:     make(chan Message),
	}
	model.printInfos()
	return model
}

func GetSettingsFile() *os.File {
	// open the file of settings
	jsonFile, err := os.Open("parametres.json")
	if err != nil {
		fmt.Println("Aucun fichier de paramètres trouvé, création d'un nouveau vide")
		makeSettingsFile(nil)
		return nil
	}
	return jsonFile
}

// Function that build up a new empty parameters file
func makeSettingsFile(settings *Settings) {
	if settings == nil {
		settings = getNewEmptySettings()
	}
	// file init
	jsonString, _ := json.MarshalIndent(settings, "", "\t")
	_ = ioutil.WriteFile("parametres.json", jsonString, 0644)
}

// Function that read the parameters
func readSettingsFromFile(jsonFile *os.File) *Settings {
	defer jsonFile.Close()

	// lire le fichier
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("%s", err)
		os.Exit(1)
	}
	// Init the struct of settings
	settings := Settings{}
	json.Unmarshal(byteValue, &settings)

	return &settings
}

//
func (model *Model) changeSettings(data interface{}) {
	fmt.Println("Changing the settings...")
	model.Settings.APIKey = (data.(map[string]interface{}))["api-key"].(string)
	model.Settings.Proxy = (data.(map[string]interface{}))["proxy"].(string)
	makeSettingsFile(model.Settings)
	sendMessage(&Message{
		Name: "test callback",
		Data: map[string]string{
			"message": "Modification paramètres : OK",
			"type":    "SUCCESS",
		},
	})
	return
}

func (model *Model) printInfos() {
	fmt.Println("Model info :")
	if model.client != nil {
		fmt.Println("\t - Client connected")
	} else {
		fmt.Println("\t - No client connected")
	}
	fmt.Println("\t - Settings :")
	fmt.Printf("\t\t . Proxy : %s\n", model.Settings.Proxy)
	fmt.Printf("\t\t . APIKey : %s\n", model.Settings.APIKey)
}

// Message Handling
func messangeHandling() {
	defer theModel.client.Close()
	fmt.Println("New websocket connected !")
	for {
		m := Message{}

		err := theModel.client.ReadJSON(&m)
		if err != nil {
			fmt.Println("Error reading json.", err)
		}

		fmt.Printf("Got message: %#v\n", m)

		switch m.Name {
		case "settings change":
			theModel.changeSettings(m.Data)
			theModel.printInfos()
		case "test request":
			go makeRequest()
		default:
			fmt.Println("DO NOTHING !")
		}

		theModel.mu.Lock()
		err = theModel.client.WriteJSON(m)
		theModel.mu.Unlock()

		if err != nil {
			fmt.Println(err)
			alreadyConnected = false
			theModel.client = nil
			fmt.Println("New connection allowed")
			return
		}
	}
}

func sendMessage(m *Message) {
	go func() { theModel.send <- *m }()
	return
}

func write() {
	for msg := range theModel.send {
		fmt.Printf("Sending : %#v \n", msg)
		theModel.mu.Lock()
		theModel.client.WriteJSON(msg)
		theModel.mu.Unlock()
	}
}
