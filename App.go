package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

var (
	ENDPOINT = "https://1fzqk3npw4.execute-api.us-east-1.amazonaws.com/nearby_store_stage/sv"
)

type DataStruc struct {
	Center struct {
		Lat float64 `json:"lat"`
		Lag float64 `json:"lag"`
	} `json:"center"`
	Zoom        int     `json:"zoom"`
	CountryCode string  `json:"country_code"`
	East        float64 `json:"east"`
	North       float64 `json:"north"`
	Sourth      float64 `json:"sourth"`
	Westh       float64 `json:"westh"`
}

func main() {
	getData()
}

func GetInformacion() {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyDkqU-_RGtijJp-UVsndAb7No99xfeJlfk"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	//r2 := &maps.FindPlaceFromTextInputTypeTextQuery("")
	r2 := &maps.PlaceAutocompleteRequest{
		Input:        "",
		Offset:       0,
		Location:     nil,
		Origin:       nil,
		Radius:       0,
		Language:     "",
		Types:        "",
		Components:   nil,
		StrictBounds: false,
		SessionToken: maps.PlaceAutocompleteSessionToken{},
	}

	algo, err := c.PlaceAutocomplete(context.Background(), r2)
	if err != nil {
		log.Fatal("Ocurrui un error")
	}
	pretty.Println(algo)
	/*
		r := &maps.DirectionsRequest{
			Origin:      "Sydney",
			Destination: "Perth",
		}

		route, _, err := c.Directions(context.Background(), r)
		if err != nil {
			log.Fatalf("fatal error: %s", err)
		}

		pretty.Println(route)

	*/
}

func getData() {
	fmt.Println("Hola mundo")

	dt := DataStruc{
		Center: struct {
			Lat float64 `json:"lat"`
			Lag float64 `json:"lag"`
		}{Lat: 13.6441687, Lag: -89.11571889999999},
		Zoom:        16,
		CountryCode: "sv",
		East:        -89.09770518425599,
		North:       13.653760482405858,
		Sourth:      13.634576527797687,
		Westh:       -89.13373261574402,
	}
	jsonData, _ := json.Marshal(dt)

	client := http.Client{}

	request, _ := http.NewRequest("POST", ENDPOINT, bytes.NewBuffer(jsonData))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Access-Control-Allow-Origin", "https://www.tiendacercasv.com")
	request.Header.Set("Content-Length", "229")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("User-Agent", "insomnia/7.1.1")
	request.Header.Set("Accept", "*/*")
	response, err := client.Do(request)

	if err != nil {
		fmt.Println(fmt.Sprintf("ocurrio un error %e", err))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}
