package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mouckatron/portableops-stationmapper/internal/models"
)

func addStation(callsign string, maidenhead string) {
	s := &models.Station{Callsign: cleanCallsign(callsign), Maidenhead: maidenhead}
	s.MaidenheadToLatLon()

	json, err := json.Marshal(*s)
	if err != nil {
		log.Println(err)
		return
	}

	for _, url := range config.stationmapperURL {
		go apiPutRequest(fmt.Sprintf("%s/stations", url), bytes.NewBuffer(json))
	}
}

func apiPutRequest(url string, data *bytes.Buffer) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(resp.StatusCode)
}

func cleanCallsign(call string) string {
	return strings.Trim(call, "<>")
}
