package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"html/template"

	"github.com/apex/log"
	"github.com/gorilla/mux"
)

var views = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	addr := ":" + os.Getenv("PORT")
	app := mux.NewRouter()
	app.HandleFunc("/", handleIndex).Methods("GET")
	if err := http.ListenAndServe(addr, app); err != nil {
		log.WithError(err).Fatal("error listening")
	}
}

func johorReading() (pasirgudang int, err error) {

	// Air Pollutant Index of Malaysia
	// http://apims.doe.gov.my/public_v2/api_table.html
	type MalaysiaAPI struct {
		Two4HourAPI [][]string `json:"24hour_api"`
	}

	resp, err := http.Get("http://apims.doe.gov.my/data/public/CAQM/last24hours.json")
	if err != nil {
		return pasirgudang, err
	}
	var aq MalaysiaAPI
	err = json.NewDecoder(resp.Body).Decode(&aq)
	if err != nil {
		return pasirgudang, err
	}
	defer resp.Body.Close()
	// log.Infof("%v", aq)

	for _, v := range aq.Two4HourAPI {
		if v[1] == "Location" {
			latest := v[len(v)-1]
			log.Infof("Latest: %s", latest)
		}
		if v[1] == "Pasir Gudang" {
			latest := v[len(v)-1]
			log.Infof("Latest: %s", latest)
			_, err := fmt.Sscanf(latest, "%d**", &pasirgudang)
			if err != nil {
				return pasirgudang, err
			}
			break
		}
	}

	return pasirgudang, err

}

func singaporeReading() (north int, err error) {
	type SingaporePM25 struct {
		RegionMetadata []struct {
			Name          string `json:"name"`
			LabelLocation struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"label_location"`
		} `json:"region_metadata"`
		Items []struct {
			Timestamp       time.Time `json:"timestamp"`
			UpdateTimestamp time.Time `json:"update_timestamp"`
			Readings        struct {
				Pm25OneHourly struct {
					West    int `json:"west"`
					East    int `json:"east"`
					Central int `json:"central"`
					South   int `json:"south"`
					North   int `json:"north"`
				} `json:"pm25_one_hourly"`
			} `json:"readings"`
		} `json:"items"`
		APIInfo struct {
			Status string `json:"status"`
		} `json:"api_info"`
	}

	resp, err := http.Get("https://api.data.gov.sg/v1/environment/pm25")
	if err != nil {
		return north, err
	}
	var aq SingaporePM25
	err = json.NewDecoder(resp.Body).Decode(&aq)
	defer resp.Body.Close()
	log.Infof("%v", aq)
	return aq.Items[0].Readings.Pm25OneHourly.North, err
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("UP_STAGE") != "production" {
		w.Header().Set("X-Robots-Tag", "none")
	}

	sg, err := singaporeReading()
	if err != nil {
		log.WithError(err).Fatal("failed to get singapore reading")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	johor, err := johorReading()
	if err != nil {
		log.WithError(err).Fatal("failed to get johor reading")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = views.ExecuteTemplate(w, "index.html", map[string]int{
		"Singapore": sg,
		"Johor":     johor,
	})

	if err != nil {
		log.WithError(err).Fatal("template failed to parse")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
