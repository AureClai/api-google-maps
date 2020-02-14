package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Structs
type Result struct {
	Destinations []string `json:"destination_addresses"`
	Origines     []string `json:"origin_addresses"`
	Statut       string   `json:"status"`
	Rows         []Row    `json:"rows"`
}

type Row struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	Distance          Entry  `json:"distance"`
	Duration          Entry  `json:"duration"`
	DistanceInTraffic Entry  `json:"duration_in_traffic"`
	Statut            string `json:"status"`
}

type Entry struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

// Request part
func makeRequest() {
	// Creating the proyx URL
	proxyURL, err := url.Parse(theModel.Settings.Proxy)
	if err != nil {
		fmt.Println("Adresse proxy invalide")
		sendMessage(&Message{
			Name: "test callback",
			Data: map[string]string{
				"type":    "FAIL",
				"message": "Erreur Proxy : Adresse invalide",
			},
		})
		return
	}

	// Creating the URL to loaded through the proxy
	urlStr := "https://maps.googleapis.com/maps/api/distancematrix/json?origins=48.8840557,2.463851&destinations=48.8715895,2.4326677&departure_time=now&traffic_model=best_guess&key=" + theModel.Settings.APIKey
	url, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println(err)
		sendMessage(&Message{
			Name: "test callback",
			Data: map[string]string{
				"type":    "FAIL",
				"message": "Erreur requête : Url invalide",
			},
		})
		return
	}

	// Adding the proxy settings to the Transport object
	transport := http.Transport{}
	if theModel.Settings.Proxy != "" {
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	// Client Initialisation
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &transport,
	}

	// Generating the HTTP, GET request
	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		fmt.Println(err)

		sendMessage(&Message{
			Name: "test callback",
			Data: "Erreur requête : Construction de la requête impossible",
		})
		return
	}

	// calling the URL
	r, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		sendMessage(&Message{
			Name: "test callback",
			Data: map[string]string{
				"type":    "FAIL",
				"message": "Erreur requête : Impossible d'effectuer la requête (proxy ?)",
			},
		})
		return
	}
	defer r.Body.Close()

	// Init Results struct
	res := new(Result)
	json.NewDecoder(r.Body).Decode(res)

	fmt.Printf("Got : %#v\n", res)

	if res.Statut != "OK" {
		sendMessage(&Message{
			Name: "test callback",
			Data: map[string]string{
				"type":    "FAIL",
				"message": "Erreur requête : " + res.Statut,
			},
		})
		return
	}

	messagePayload := fmt.Sprintf("Requête OK entre %v et %v : %v sec", res.Origines[0], res.Destinations[0], res.Rows[0].Elements[0].DistanceInTraffic.Value)
	go AddRecord(res)

	// Sending response
	sendMessage(&Message{
		Name: "test callback",
		Data: map[string]string{
			"type":    "SUCCESS",
			"message": messagePayload,
		},
	})

	return
}
