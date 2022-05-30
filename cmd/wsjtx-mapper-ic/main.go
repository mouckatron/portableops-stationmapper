package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/k0swe/wsjtx-go/v4"
	"github.com/mouckatron/portableops-stationmapper/internal/models"
)

var (
	maidenheadAll = regexp.MustCompile("([A-R]{2}[0-9]{2}([a-x]{2}([0-9]{2})?)?)$")
	maidenhead4   = regexp.MustCompile(".*([A-R]{2}[0-9]{2})$")
	maidenhead6   = regexp.MustCompile("([A-R]{2}[0-9]{2}[a-x]{2})$")
	maidenhead8   = regexp.MustCompile("([A-R]{2}[0-9]{2}[a-x]{2}[0-9]{2})$")
	responseMsg   = regexp.MustCompile("(R?[-+]?[0-9]+|RRR|(RR)?73)$")
)

func main() {

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

func parseDecodeMessage(message wsjtx.DecodeMessage) {
	// toJSON, err := json.MarshalIndent(message, "", "  ")

	// if err != nil {
	// 	log.Println("Error:", err)
	// 	return
	// }

	// log.Println(string(toJSON))

	switch true {
	case strings.HasPrefix(message.Message, "CQ "):
		parseDecodedCQ(message)
	default:
		parseDecodedToFrom(message)
	}
}

func parseDecodedCQ(message wsjtx.DecodeMessage) {
	parts := strings.Fields(message.Message)
	maidenhead := ""

	callsign := parts[1]
	if len(parts) == 4 { // to match something like CQ DX
		callsign = parts[2]
	}

	if len(parts) > 2 {
		maidenheadMatch := maidenheadAll.FindStringSubmatch(message.Message)
		if maidenheadMatch != nil {
			maidenhead = maidenheadMatch[0]
		}
	}

	log.Println("Callsign:", callsign, "Location:", maidenhead)
	go addStation(callsign, maidenhead)
}

func parseDecodedToFrom(message wsjtx.DecodeMessage) {
	parts := strings.Fields(message.Message)
	if len(parts) < 2 {
		return
	}

	to, from, maidenhead := "", "", ""

	if responseMsg.MatchString(parts[len(parts)-1]) {
		to = parts[0]
		from = parts[1]
	} else if maidenheadAll.MatchString(message.Message) {
		to = parts[0]
		from = parts[1]
		maidenhead = parts[2]
	}
	log.Println("To:", to, "From:", from, "Location:", maidenhead)
	go addStation(from, maidenhead)
}

func addStation(callsign string, maidenhead string) {
	s := &models.Station{Callsign: cleanCallsign(callsign), Maidenhead: maidenhead}
	s.MaidenheadToLatLon()

	client := &http.Client{}

	json, err := json.Marshal(*s)
	if err != nil {
		log.Println(err)
		return
	}

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/stations", bytes.NewBuffer(json))
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
