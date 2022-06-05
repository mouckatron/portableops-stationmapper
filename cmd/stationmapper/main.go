package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mouckatron/portableops-stationmapper/cmd/stationmapper/ui"
)

func main() {

	var myHost string
	const myHostHelp = "Address stationmapper will be available at"
	var myPort int
	const myPortHelp = "Port stationmapper will be available at"
	var ginProduction = true

	// normal arguments
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.StringVar(&myHost, "host", "0.0.0.0", myHostHelp)
	flag.IntVar(&myPort, "port", 8080, myPortHelp)

	flag.Parse()

	cmd := flag.Arg(0)

	readConfig()
	setupDB()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	if cmd == "serve" {
		serveHTTP(myHost, myPort, ginProduction)
	} else if cmd == "tiles" {
		GetTiles()
	}

}

func cleanup() {
	cleanupDB()
}

func serveHTTP(myHost string, myPort int, ginProduction bool) {
	// start the API server which uses the RigConnection thing
	if ginProduction {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(cors.Default())
	tileRouterPaths(router)
	// routerPaths(router)
	ui.RouterPaths(router)
	APIRouterPaths(router)
	router.Run(fmt.Sprintf("%s:%d", myHost, myPort))

}
