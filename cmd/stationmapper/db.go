package main

import (
	"log"

	"github.com/mouckatron/portableops-stationmapper/internal/models"
	//"gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func setupDB() {
	var err error = nil

	db, err = gorm.Open(sqlite.Open("stationmapper.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&DBStation{})
}

func cleanupDB() {

}

func DBUpsertStation(station models.Station) {
	var dbStation DBStation

	result := db.Where("Callsign = ?", station.Callsign).Find(&dbStation)

	if result.RowsAffected > 0 {
		// some more complicated update
		if station.Maidenhead != "" {
			station.MaidenheadToLatLon()
			db.Model(&dbStation).Update("Maidenhead", station.Maidenhead)
			db.Model(&dbStation).Update("Lat", station.Lat)
			db.Model(&dbStation).Update("Lon", station.Lon)
		}

	} else {
		log.Println("Creating new station", station.Callsign)
		dbStation.Callsign = station.Callsign
		dbStation.Maidenhead = station.Maidenhead
		station.MaidenheadToLatLon()
		dbStation.Lat = station.Lat
		dbStation.Lon = station.Lon
		db.Create(&dbStation)
	}

}

type DBStation struct {
	gorm.Model
	Callsign   string `gorm:"unique"`
	Maidenhead string
	Lat        float64
	Lon        float64
}
