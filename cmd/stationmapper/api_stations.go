package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mouckatron/portableops-stationmapper/internal/models"
)

//var stations = make(map[string]*models.Station) // TODO: Get from DB

func apiGetStations(c *gin.Context) {

	var dbStations []DBStation
	db.Find(&dbStations)

	var response = make(map[string]*models.Station)

	for _, s := range dbStations {
		response[s.Callsign] = &models.Station{Callsign: s.Callsign,
			Maidenhead: s.Maidenhead,
			Lat:        s.Lat,
			Lon:        s.Lon}

		if s.Lat == 0 || s.Lon == 0 {
			response[s.Callsign].MaidenheadToLatLon()
		}

	}
	c.JSON(200, response)
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

	// if station, ok := stations[incoming.Callsign]; ok {
	// 	station.Merge(incoming)

	// } else {
	// 	stations[incoming.Callsign] = &incoming
	DBUpsertStation(incoming)
	// }
}
