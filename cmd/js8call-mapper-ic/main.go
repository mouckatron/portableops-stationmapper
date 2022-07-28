package main

import (
	"log"
  "reflect"

	"github.com/mouckatron/js8call-go"
)

func main() {

	loadConfig()

  client := js8call.MakeTCPClient("localhost", "2442")
  client.Debug = true
	go client.Start()

	for {
		select {
		case message := <-client.MessageChannel:
			go handleMessage(message)
		case error := <-client.ErrorChannel:
			log.Println(error)
		}
	}
}

func handleMessage(message interface{}) {
  switch message.(type) {
  case js8call.RXSpot:
    _message := message.(js8call.RXSpot)
    if _message.Params.GRID != "" {
      go addStation(_message.Params.CALL, _message.Params.GRID)
    }
  default:
    log.Println("Other:", reflect.TypeOf(message), message)
  }
}
