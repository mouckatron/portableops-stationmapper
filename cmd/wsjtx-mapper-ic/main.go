package main

import (
	"log"
	"reflect"
	"regexp"

	"github.com/k0swe/wsjtx-go/v4"
)

var (
	maidenheadAll = regexp.MustCompile("([A-R]{2}[0-9]{2}([a-x]{2}([0-9]{2})?)?)$")
	maidenhead4   = regexp.MustCompile(".*([A-R]{2}[0-9]{2})$")
	maidenhead6   = regexp.MustCompile("([A-R]{2}[0-9]{2}[a-x]{2})$")
	maidenhead8   = regexp.MustCompile("([A-R]{2}[0-9]{2}[a-x]{2}[0-9]{2})$")
	responseMsg   = regexp.MustCompile("(R?[-+][0-9]+|RRR|(RR)?73)$")
)

func main() {

	loadConfig()

	wsjtxServer, err := wsjtx.MakeServer()
	if err != nil {
		log.Fatalf("%v", err)
	}
	wsjtxChannel := make(chan interface{}, 5)
	errChannel := make(chan error, 5)
	go wsjtxServer.ListenToWsjtx(wsjtxChannel, errChannel)

	for {
		select {
		case message := <-wsjtxChannel:
			go handleServerMessage(message)
		case error := <-errChannel:
			log.Println(error)
		}
	}

}

// When we receive WSJT-X messages, display them.
func handleServerMessage(message interface{}) {
	switch message.(type) {
	case wsjtx.HeartbeatMessage:
		// log.Println("Heartbeat:", message)
	case wsjtx.StatusMessage:
	// 	log.Println("Status:", message)
	case wsjtx.DecodeMessage:
		log.Println("Decode:", message)
		parseDecodeMessage(message.(wsjtx.DecodeMessage))
	case wsjtx.ClearMessage:
		// log.Println("Clear:", message)
	case wsjtx.QsoLoggedMessage:
		// log.Println("QSO Logged:", message)
	case wsjtx.CloseMessage:
		// log.Println("Close:", message)
	case wsjtx.WSPRDecodeMessage:
		// log.Println("WSPR Decode:", message)
	case wsjtx.LoggedAdifMessage:
		// log.Println("Logged Adif:", message)
	default:
		log.Println("Other:", reflect.TypeOf(message), message)
	}
}
