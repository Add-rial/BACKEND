package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const URL = "https://api.wheretheiss.at/v1/satellites/25544"

func main() {
	body := get_request()
	latitude, longitude := decodeJSON(body)
	fmt.Println("Latitude: ", latitude)
	fmt.Println("Longitude: ", longitude)
}

func get_request() []byte {
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("Error encountered")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error encountered")
	}
	return body
}

func decodeJSON(dataJSON []byte) (latitude float64, longitude float64) {
	var dataDecoded map[string]interface{}

	json.Unmarshal(dataJSON, &dataDecoded)
	latitude, _ = dataDecoded["latitude"].(float64)
	longitude, _ = dataDecoded["longitude"].(float64)
	return
}
