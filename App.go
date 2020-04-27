package main

import (
	"CoordenadasTiendaCercaSV/modelo"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"googlemaps.github.io/maps"
)

var (
	ENDPOINT  = "https://1fzqk3npw4.execute-api.us-east-1.amazonaws.com/nearby_store_stage/sv"
	CLAVEJUAN = "AIzaSyCswu8im_YgIcNBGFmRr-gRVBLqBHwXVxk"
	//CLAVETIENDA="AIzaSyDkqU-_RGtijJp-UVsndAb7No99xfeJlfk"
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
	GetInformacion()
}

func GetInformacion() {
	c, err := maps.NewClient(maps.WithAPIKey(CLAVEJUAN))
	if err != nil {
		log.Fatalf("fatal error: %s", err.Error())
	}

	//r2 := &maps.FindPlaceFromTextInputTypeTextQuery("")
	r3 := &maps.PlaceDetailsRequest{
		PlaceID:      "ChIJL13XQfeXYo8ROb7a0b7clBw",
		Language:     "",
		Fields:       nil,
		SessionToken: maps.PlaceAutocompleteSessionToken{},
		Region:       "",
	}

	algo2, err := c.PlaceDetails(context.Background(), r3)
	if err != nil {
		log.Fatal(fmt.Sprintf("Ocurrio un error  %e", err))
	}

	//center
	Centerlat := algo2.Geometry.Location.Lat
	Centerlag := algo2.Geometry.Location.Lng

	//ubicaciones de
	north := algo2.Geometry.Viewport.NorthEast.Lat
	east := algo2.Geometry.Viewport.NorthEast.Lng
	south := algo2.Geometry.Viewport.SouthWest.Lat
	west := algo2.Geometry.Viewport.SouthWest.Lng

	dt := DataStruc{
		Center: struct {
			Lat float64 `json:"lat"`
			Lag float64 `json:"lag"`
		}{Lat: Centerlat, Lag: Centerlag},
		Zoom:        16,
		CountryCode: "sv",
		East:        east,
		North:       north,
		Sourth:      south,
		Westh:       west,
	}
	getData(dt)

}

func getData(dt DataStruc) {
	fmt.Println("Hola mundo")

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

func GetDepartamentos() modelo.Departamentos {
	ManejadorDepartamentos, err := ioutil.ReadFile("Departamentos.json")
	if err != nil {
		log.Fatal(fmt.Sprintf("No se pudo leer el archivo %e", err))
	}
	dp := modelo.Departamentos{}
	errJson := json.Unmarshal(ManejadorDepartamentos, &dp)
	if errJson != nil {
		log.Fatal(fmt.Sprintf("Ocurrio un error al parsear el archivo %e", err))
	}
	return dp
}
