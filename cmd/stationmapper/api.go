package main

import (
	"github.com/gin-gonic/gin"
)

func APIRouterPaths(r *gin.Engine) {
	r.GET("/stations", apiGetStations)
	r.PUT("/stations", apiPutStation)

}
