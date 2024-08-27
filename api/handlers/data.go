package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type PageData struct {
	Resultat         string
	Title            string
	ArtistItem       []Artist
	DateLocationItem DateLocationsData
	DatesItem        DatesData
	LocationItem     LocationsData
}

type Err_Data struct {
	Title   string
	Content string
}

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	LocationDate map[string][]string
	NextId       int
	PreviousId   int
}

type DateLocation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type DateLocationsData struct {
	Index []DateLocation `json:"index"`
}

type Date struct {
	ID   int      `json:"id"`
	Date []string `json:"dates"`
}

type DatesData struct {
	Index []Date `json:"index"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type LocationsData struct {
	Index []Location `json:"index"`
}

func FetchData(url string, inter interface{}) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	_ = json.Unmarshal(responseData, &inter)
}

func MergeData(data PageData) {
	for i := 0; i < len(data.ArtistItem); i++ {

		if data.ArtistItem[i].LocationDate == nil {
			data.ArtistItem[i].LocationDate = make(map[string][]string)
		}
		for key, value := range data.DateLocationItem.Index[i].DatesLocations {
			data.ArtistItem[i].LocationDate[key] = value
		}

	}
}
