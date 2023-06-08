package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
)

type Astronaut struct {
	Name  string `json:"name"`
	Craft string `json:"craft"`
}

type AstronautsResponse struct {
	People []Astronaut `json:"people"`
}

func main() {
	// Fetch data from the API
	response, err := http.Get("http://api.open-notify.org/astros.json")
	if err != nil {log.Fatal("Failed to fetch data from API:", err)}
	defer response.Body.Close()

	// Parse the JSON response
	var data AstronautsResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {log.Fatal("Failed to decode JSON response:", err)}

	// Sort the astronauts by craft and then name
	sort.Slice(data.People, func(i, j int) bool {
		if data.People[i].Craft == data.People[j].Craft {
			return data.People[i].Name < data.People[j].Name
		}
		return data.People[i].Craft < data.People[j].Craft
	})

	// Create a CSV file
	file, err := os.Create("astronauts.csv")
	if err != nil {log.Fatal("Failed to create file:", err)}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"Name", "Craft"}
	err = writer.Write(header)
	if err != nil {log.Fatal("Failed to write CSV header:", err)}

	// Write each astronaut to the CSV file
	for _, astronaut := range data.People {
		row := []string{astronaut.Name, astronaut.Craft}
		err = writer.Write(row)
		if err != nil {log.Fatal("Failed to write CSV row:", err)}
	}

	log.Println("CSV file created successfully!")
}
