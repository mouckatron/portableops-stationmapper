package main

import (
	"fmt"
	"log"

	"github.com/mouckatron/portableops-stationmapper/internal/models"
	"github.com/spf13/viper"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func setupDB() {
	var err error = nil

	switch viper.GetString("db.type") {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(viper.GetString("db.path")), &gorm.Config{})
	case "postgres":
		dsn := fmt.Sprintf("user=stationmapper dbname=stationmapper")
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.DBStation{})
}

func cleanupDB() {

}

func DBUpsertStation(station models.Station) {
	var dbStation models.DBStation

	result := db.Where("Callsign = ?", station.Callsign).Find(&dbStation)

	if result.RowsAffected > 0 {
		// some more complicated update
		if station.Maidenhead != "" {
			log.Println("Updating station", station.Callsign)
			station.MaidenheadToLatLon()

			dbStation.Maidenhead = station.Maidenhead
			dbStation.Lat = station.Lat
			dbStation.Lon = station.Lon
			db.Save(&dbStation)
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
