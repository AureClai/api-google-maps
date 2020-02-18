package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"googlemaps.github.io/maps"
)

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

	// Generating the request client
	c, err := maps.NewClient(maps.WithAPIKey(theModel.Settings.APIKey), maps.WithHTTPClient(client))
	if err != nil {
		fmt.Println(err)
		sendMessage(&Message{
			Name: "test callback",
			Data: map[string]string{
				"type":    "FAIL",
				"message": "Erreur requête : Construction du client Maps impossible",
			},
		})
	}

	r := &maps.DistanceMatrixRequest{
		Origins:       []string{"48.8840557,2.463851"},
		Destinations:  []string{"48.8715895,2.4326677"},
		DepartureTime: "now",
		Mode:          maps.TravelModeDriving,
		Units:         maps.UnitsMetric,
		TrafficModel:  maps.TrafficModelBestGuess,
	}

	response, err := c.DistanceMatrix(context.Background(), r)
	if err != nil {
		fmt.Println(err)

		sendMessage(&Message{
			Name: "test callback",
			Data: map[string]string{
				"type":    "FAIL",
				"message": "Erreur requête : Impossible d'effectuer la requête (proxy ?)",
			},
		})
	}
	fmt.Printf("Got : %#v\n", response)

	messagePayload := fmt.Sprintf("Requête OK entre %v et %v : %v sec", response.OriginAddresses[0], response.DestinationAddresses[0], response.Rows[0].Elements[0].DurationInTraffic.Seconds())
	go AddRecord(response)

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
