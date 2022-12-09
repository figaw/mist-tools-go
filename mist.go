package mist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"io"
)

type envelope struct {
	MessageId string `json:"messageId"`
	TraceId   string `json: traceId`
	Payload   string `json: payload`
}

type handlers map[string]func(string)
type iFunc func()

func MistService(hs handlers) {
	action := os.Args[len(os.Args)-2]
	handler := hs[action]
	if handler != nil {
		invokeHandler(action, handler)
	}
}

func MistServiceWithInit(hs handlers, init iFunc) {
	action := os.Args[len(os.Args)-2]
	handler := hs[action]
	if handler != nil {
		invokeHandler(action, handler)
	} else if init != nil {
		fmt.Println("running init")
		init()
	}
}

func invokeHandler(action string, handler func(string)) {
	fmt.Printf("running handler for %s\n", action)
	var e envelope
	err := json.Unmarshal([]byte(os.Args[len(os.Args)-1]), &e)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("contents of decoded json is: %#v\r\n", e)
	fmt.Println(e.Payload)

	handler(e.Payload)
}

func PostToRapid[T interface{}](event string, reply T) {
	fmt.Printf("contents of reply is: %#v\r\n", reply)

	body, _ := json.Marshal(reply)
	fmt.Printf("posting %s to (%s/%s)\n", string(body), os.Getenv("RAPID"), event)

	PostBodyToRapid(event, bytes.NewBuffer(body))
}

func PostBodyToRapid(event string, body io.Reader) {
	resp, err := http.Post(fmt.Sprintf("%s/%s", os.Getenv("RAPID"), event), "application/json", body)
	if err != nil {
		fmt.Println("Get failed with error: ", err)
	}
	defer resp.Body.Close()
}
