package main

import (
	"CoordenadasTiendaCercaSV/modelo"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
	"gopkg.in/h2non/gentleman.v2/plugins/body"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	//"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo/v4"
	"gopkg.in/h2non/gentleman.v2"
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

	serverRun()
	//ID :=GetIDMunicipio("Santiago Texacuangos , San Salvador")
	//pretty.Println(ID)
	//GetInformacion()
}

//FindPlaceByID Busca un lugar por ID
func FindPlaceByID(ID string) maps.PlaceDetailsResult {
	c, err := maps.NewClient(maps.WithAPIKey(CLAVEJUAN))
	if err != nil {
		log.Fatalf("fatal error: %s", err.Error())
	}
	r3 := &maps.PlaceDetailsRequest{
		PlaceID: ID,
		//"ChIJCVt4J6lPY48R77BngjCT2C8",
		//"ChIJL13XQfeXYo8ROb7a0b7clBw",
		Language:     "",
		Fields:       nil,
		SessionToken: maps.PlaceAutocompleteSessionToken{},
		Region:       "",
	}

	DetalleLugar, err := c.PlaceDetails(context.Background(), r3)
	if err != nil {
		log.Fatal(fmt.Sprintf("Ocurrio un error  %e", err))
	}
	//pretty.Println(DetalleLugar)
	return DetalleLugar
}
func GetInformacion() {

	//r2 := &maps.FindPlaceFromTextInputTypeTextQuery("")

	fmt.Println("--------------------------------------------------------------------")
	departamentos := GetDepartamentos()
	//pretty.Println(departamentos)

	for _, departamento := range departamentos.ListaDepartamentos {
		//pretty.Println(departamento.Nombre)
		for _, municipio := range departamento.Municipios {
			//fmt.Println(fmt.Sprintf("Departamento: %s - Municipio: %s",departamento.Nombre,municipio))
			IDMunicipio := GetIDMunicipio(municipio + "," + departamento.Nombre)
			//pretty.Println(fmt.Sprintf("Departamento: %s  - Municipio: %s  IDMunicipio : %s",departamento.Nombre,municipio,IDMunicipio))
			//detalleLugar := FindPlaceByID(IDMunicipio)

			FindDatosTiendaCerca(IDMunicipio)
		}
	}

	//getData(dt)

}

func FindDatosTiendaCerca(ID string) {
	DetalleLugar := FindPlaceByID(ID)
	//center
	Centerlat := DetalleLugar.Geometry.Location.Lat
	Centerlag := DetalleLugar.Geometry.Location.Lng

	//ubicaciones de
	north := DetalleLugar.Geometry.Viewport.NorthEast.Lat
	east := DetalleLugar.Geometry.Viewport.NorthEast.Lng
	south := DetalleLugar.Geometry.Viewport.SouthWest.Lat
	west := DetalleLugar.Geometry.Viewport.SouthWest.Lng

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
	getDatosTienda(dt)
	//pretty.Println(dt)
	//getData(dt)
}
func getDatosTienda(dt DataStruc) {
	/*
		jsonData ,err  :=json.Marshal(dt)

		if err != nil {
			log.Fatal(fmt.Sprintf("Ocurrio un error %e",err))
		}
	*/
	client := gentleman.New()
	client.URL(ENDPOINT)

	req := client.Request()

	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Access-Control-Allow-Origin", "https://www.tiendacercasv.com")
	req.Method("POST")
	req.Use(body.JSON(dt))

	res, errorRequest := req.Send()
	if errorRequest != nil {
		log.Fatal(fmt.Sprintf("ocurrion un error %e", errorRequest))
	}
	pretty.Println(res.String())
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

//Obtine los de departamentos desde archivo json
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

//Busca El ID de cada Departamento segun su nombre
func GetIDMunicipio(lugar string) string {
	//datos := &maps.place
	c, err := maps.NewClient(maps.WithAPIKey(CLAVEJUAN))
	if err != nil {
		log.Fatal(fmt.Sprintf("Ocurrio un error %e", err))
	}
	findLugar := &maps.FindPlaceFromTextRequest{
		Input:                 lugar,
		InputType:             "textquery",
		Fields:                nil,
		LocationBias:          "",
		LocationBiasPoint:     nil,
		LocationBiasCenter:    nil,
		LocationBiasRadius:    0,
		LocationBiasSouthWest: nil,
		LocationBiasNorthEast: nil,
	}

	result, errorFindPlace := c.FindPlaceFromText(context.Background(), findLugar)

	if errorFindPlace != nil {
		log.Fatal(fmt.Sprintf("error findplace %e", errorFindPlace))
	}

	pretty.Println(result)
	var PlaceID string
	for _, can := range result.Candidates {
		PlaceID = can.PlaceID
	}
	return PlaceID
}

func serverRun() {
	e := echo.New()
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	e.GET("/api/v1/departamentos/:municipio", func(c echo.Context) error {

		municiopio := c.Param("municipio")

		IDMuniciopio := GetIDMunicipio(municiopio)

		datos := FindPlaceByID(IDMuniciopio)

		return c.JSON(http.StatusOK, datos)
	})
	e.GET("/api/v1/departamentos", func(c echo.Context) error {
		departamentos := GetDepartamentos()
		return c.JSON(http.StatusOK, departamentos)
	})
	e.Logger.Fatal(e.Start(":" + port))
}
