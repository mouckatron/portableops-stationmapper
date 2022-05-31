package main

import (
	"log"
	"strings"

	"github.com/k0swe/wsjtx-go/v4"
)

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
	if len(parts) < 3 {
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
