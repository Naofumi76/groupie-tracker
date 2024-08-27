package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var (
	index, _         = template.ParseFiles("web/templates/index.html")
	artTemplate, _   = template.ParseFiles("web/templates/artist.html")
	errorTemplate, _ = template.ParseFiles("web/templates/error.html")
	data             = PageData{
		Resultat: "",
		Title:    "Groupie-Tracker",
	}
	err_data = Err_Data{
		Title:   "404",
		Content: "Page not found",
	}
	art = Artist{}
)

// Handler of the main page
func HandleHome(w http.ResponseWriter, r *http.Request) {
	MergeData(data)

	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if len(data.ArtistItem) == 0 {
		err_data.Title = "500"
		err_data.Content = "Internal server error"
		errorTemplate.Execute(w, err_data)
		log.Printf("Error 500: unreachable data")
		err_data.Title = "404"
		err_data.Content = "Page not found"
		return
	} else if r.URL.Path == "/" {
		membersStr := r.URL.Query()["members"]
		creationDateStrMin := r.FormValue("creation-date-min")
		creationDateStrMax := r.FormValue("creation-date-max")
		locationStr := r.FormValue("location")
		firstAlbumYearStrMin := r.FormValue("first-album-year-min")
		firstAlbumYearStrMax := r.FormValue("first-album-year-max")

		var filteredItems []Artist

		// Parse creation date filter
		creationDateFilterMin := 0
		creationDateFilterMax := 2025
		noCreationDateFilter := false
		if creationDateStrMin != "" || creationDateStrMax != "" {
			var err, err2 error
			creationDateFilterMin, err = strconv.Atoi(creationDateStrMin)
			creationDateFilterMax, err2 = strconv.Atoi(creationDateStrMax)
			if err != nil || err2 != nil {
				log.Println("Invalid creation date:", err)
				err_data.Title = "500"
				err_data.Content = "Internal server error"
				errorTemplate.Execute(w, err_data)
				return
			}
		} else {
			// If creation-date is not provided, treat it as no filter
			noCreationDateFilter = true
		}

		// Parse first album year filter
		noFirstAlbumYearFilter := false
		firstAlbumYearFilterMin := 0
		firstAlbumYearFilterMax := 2025
		if firstAlbumYearStrMin != "" || firstAlbumYearStrMax != "" {
			var err, err2 error
			firstAlbumYearFilterMin, err = strconv.Atoi(firstAlbumYearStrMin)
			firstAlbumYearFilterMax, err2 = strconv.Atoi(firstAlbumYearStrMax)
			if err != nil || err2 != nil {
				log.Println("Invalid first album year:", err)
				err_data.Title = "500"
				err_data.Content = "Internal server error"
				errorTemplate.Execute(w, err_data)
				return
			}
		} else {
			noFirstAlbumYearFilter = true
		}

		// Check if any filter is applied
		if len(membersStr) > 0 || !noCreationDateFilter || locationStr != "" || !noFirstAlbumYearFilter {
			membersMap := make(map[int]bool)
			for _, m := range membersStr {
				num, err := strconv.Atoi(m)
				if err == nil {
					membersMap[num] = true
				}
			}
			for _, artist := range data.ArtistItem {
				locationMatch := false
				for location := range artist.LocationDate {
					if location == ReverseLocStyle(strings.ToLower(locationStr)) {
						locationMatch = true
						break
					}
				}

				// Parse the first album date to get the year
				firstAlbumDate, err := time.Parse("02-01-2006", artist.FirstAlbum)
				if err != nil {
					log.Printf("Invalid first album date format for artist %s: %v", artist.Name, err)
					continue // Skip invalid entries
				}
				artistFirstAlbumYear := firstAlbumDate.Year()

				// Apply the filters
				if (len(membersMap) == 0 || membersMap[len(artist.Members)]) &&
					(noCreationDateFilter || artist.CreationDate >= creationDateFilterMin) && (artist.CreationDate <= creationDateFilterMax) &&
					(locationStr == "" || locationMatch) &&
					(noFirstAlbumYearFilter || (artistFirstAlbumYear >= firstAlbumYearFilterMin) && (artistFirstAlbumYear <= firstAlbumYearFilterMax)) {
					filteredItems = append(filteredItems, artist)
				}
			}
		} else {
			// If no filters are applied, show all artists
			filteredItems = data.ArtistItem
		}

		err := index.Execute(w, PageData{ArtistItem: filteredItems, Title: data.Title, LocationItem: data.LocationItem})
		if err != nil {
			log.Printf("Template execution error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
		} else {
			log.Printf("Status OK : %v", http.StatusOK)
		}
		return
	} else if strings.HasPrefix(r.URL.Path, "/artist") {
		artistIDStr := r.URL.Path[len("/artist/"):]
		artistID, err := strconv.Atoi(artistIDStr)
		if err != nil {
			log.Println("Invalid artist ID:", err)
			err_data.Title = "500"
			err_data.Content = "Internal server error"
			errorTemplate.Execute(w, err_data)
			err_data.Title = "404"
			err_data.Content = "Page not found"
			return
		}

		if artistID >= 1 && artistID <= len(data.ArtistItem) {
			art = data.ArtistItem[artistID-1]
			if artistID == 1 {
				art.NextId = art.Id + 1
				artTemplate.Execute(w, art)
				return
			}
			if artistID == len(data.ArtistItem) {
				art.PreviousId = art.Id - 1
				artTemplate.Execute(w, art)
				return
			}
			art.NextId = art.Id + 1
			art.PreviousId = art.Id - 1
			artTemplate.Execute(w, art)
			return
		} else {
			err_data.Title = "500"
			err_data.Content = "Internal server error"
			errorTemplate.Execute(w, err_data)
			err_data.Title = "404"
			err_data.Content = "Page not found"
			log.Printf("Error 404: Page not found for path %s", r.URL.Path)
			return
		}
	} else {
		errorTemplate.Execute(w, err_data)
		log.Printf("Error 404: Page not found for path %s", r.URL.Path)
		return
	}
}

func Init() {
	FetchData("https://groupietrackers.herokuapp.com/api/artists", &data.ArtistItem)
	FetchData("https://groupietrackers.herokuapp.com/api/relation", &data.DateLocationItem)
	// Fetch locations data and set it into data.LocationItem
	FetchData("https://groupietrackers.herokuapp.com/api/locations", &data.LocationItem)

	for i, loc := range data.LocationItem.Index {
		t := []string{}
		for _, place := range loc.Locations {
			if !In(strings.Title(LocStyle(place)), t) {
				t = append(t, strings.Title(LocStyle(place)))
			}
		}
		data.LocationItem.Index[i].Locations = t
	}
}

func LocStyle(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Replace(s, "-", ", ", 1)
	return s
}

func ReverseLocStyle(s string) string {
	s = strings.Replace(s, ", ", "-", 1)
	s = strings.Replace(s, " ", "_", -1)
	return s
}

func In(s string, t []string) bool {
	for _, val := range t {
		if s == val {
			return true
		}
	}
	return false
}
