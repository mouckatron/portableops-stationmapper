package models

import (
	"log"

	"github.com/pd0mz/go-maidenhead"
)

type Station struct {
	Callsign   string  `json:"callsign"`
	Maidenhead string  `json:"maidenhead"`
	Lat        float64 `json:"latitude"`
	Lon        float64 `json:"longitude"`
}

func (s *Station) Merge(incoming Station) {

	if s.Maidenhead != "" {
		s.Maidenhead = incoming.Maidenhead
	}

	if s.Lat == 0 || s.Lon == 0 {
		s.MaidenheadToLatLon()
	}
}

func (s *Station) MaidenheadToLatLon() {
	if s.Maidenhead != "" {
		p, err := maidenhead.ParseLocator(s.Maidenhead)
		if err != nil {
			log.Println(err)
			return
		}

		s.Lat = p.Latitude
		s.Lon = p.Longitude
	}
}
