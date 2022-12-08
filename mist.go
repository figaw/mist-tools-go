package mist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Envelope struct {
	messageId string  `json: messageId`
	traceId   string  `json: traceId`
	payload   Payload `json: payload`
}

type Handlers map[string]interface{}

type iFunc func()

func mistService(handlers Handlers) {
	action := os.Args[len(os.Args)-2]
	handler := handlers[action]
	if handler != nil {
		fmt.Println("running handler for %s", action)
		var envelope Envelope
		json.Unmarshal([]byte(os.Args[len(os.Args)-1]), &envelope)
		handler.(func(Payload))(envelope.payload)
	}
}

func mistServiceWithInit(handlers Handlers, init iFunc) {
	action := os.Args[len(os.Args)-2]
	handler := handlers[action]
	if handler != nil {
		fmt.Println("running handler for %s", action)
		var envelope Envelope
		json.Unmarshal([]byte(os.Args[len(os.Args)-1]), &envelope)
		handler.(func(Payload))(envelope.payload)
	} else if init != nil {
		fmt.Println("running init")
		init()
	}
}

func postToRapid(event string, reply Reply) {
	body, _ := json.Marshal(reply)
	fmt.Println("posting %s to (%s/%s)", body, os.Getenv("RAPID"), event)
	resp, err := http.Post(fmt.Sprintf("%s/%s", os.Getenv("RAPID"), event), "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Get failed with error: ", err)
	}
	defer resp.Body.Close()
}
