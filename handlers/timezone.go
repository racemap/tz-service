package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	timezone "github.com/evanoberholster/timezoneLookup"
	"github.com/sirupsen/logrus"
)

type TimezoneInfos struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	Offset int    `json:"offset"`
}

func TimezoneHandler(tz timezone.TimezoneInterface, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		latS := r.FormValue("lat")
		lngS := r.FormValue("lng")

		lat, lng, err := prepareCoordinates(latS, lngS)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Input coordinates have wrong format. Plaese deliver as float number. For example /api?lng=52.517932&lat=13.402992")
			return
		}

		tzInfos, err := getTimezoneInfos(lat, lng, tz)

		if err != nil {
			log.Error("Failed to find timezone for a request:", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to find timezone for your request.")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(tzInfos); err != nil {
			log.Error("Failed to build valid response:", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to build valid response.")
			return
		}
	}
}

func prepareCoordinates(latS string, lngS string) (float32, float32, error) {
	latF, err := strconv.ParseFloat(latS, 32)

	if err != nil {
		return -1, -1, err
	}

	lngF, err := strconv.ParseFloat(lngS, 32)

	if err != nil {
		return -1, -1, err
	}

	return float32(latF), float32(lngF), nil
}

func getTimezoneInfos(lat float32, lng float32, tz timezone.TimezoneInterface) (*TimezoneInfos, error) {
	tzName, err := tz.Query(timezone.Coord{
		Lon: lng, Lat: lat,
	})

	if err != nil {
		return nil, err
	}

	location, err := time.LoadLocation(tzName)

	if err != nil {
		return nil, err
	}

	tzID, tzOffset := time.Now().In(location).Zone()
	tzInfos := TimezoneInfos{
		Name:   tzName,
		ID:     tzID,
		Offset: tzOffset,
	}

	return &tzInfos, nil
}
