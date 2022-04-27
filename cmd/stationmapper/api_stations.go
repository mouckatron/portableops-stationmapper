package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mouckatron/portableops-stationmapper/internal/models"
	maidenhead "github.com/pd0mz/go-maidenhead"
)

var stations = make(map[string]*models.Station) // TODO: Get from DB

func apiGetStations(c *gin.Context) {

	for _, s := range stations {
		p, err := maidenhead.ParseLocator(s.Maidenhead)
		if err != nil {
			continue
		}

		s.Lat = p.Latitude
		s.Lon = p.Longitude
	}

	c.JSON(200, stations)
}

func apiPutStation(c *gin.Context) {

	var incoming models.Station

	if err := c.BindJSON(&incoming); err != nil {
		log.Println(err)
	}

	if incoming.Callsign == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	if station, ok := stations[incoming.Callsign]; ok {
		station.Merge(incoming)

	} else {
		stations[incoming.Callsign] = &incoming
	}
}

func addStation(callsign string, maidenhead string) {
	stations[callsign] = &models.Station{Callsign: callsign, Maidenhead: maidenhead}
}
