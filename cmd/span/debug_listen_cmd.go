package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type wsCmd struct {
	URLString string `long:"url" description:"Websocket endpoint" required:"yes"`
	Pretty    bool   `long:"pretty" description:"pretty-print data"`
}

func init() {
	parser.AddCommand(
		"ws",
		"Websocket listen",
		"Websocket listen debug utility",
		&wsCmd{},
	)
}

func (r *wsCmd) Execute([]string) error {
	header := http.Header{}
	header.Add("X-API-Token", opt.Token)
	dialer := websocket.Dialer{}

	conn, response, err := dialer.Dial(r.URLString, header)
	if opt.Debug {
		log.Printf("HTTP %s", response.Status)
		for k, v := range response.Header {
			log.Printf("    %s : %s", k, v)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("error reading body: %v", err)
		}

		log.Printf("BODY: '%s'", body)

	}
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.WriteMessage(1, []byte("hello there"))
	if err != nil {
		log.Fatalf("failed to send: %v", err)
	}

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("error reading message: %v", err)
		}

		log.Printf("[MSG] type=%d(%s), msg='%s'", msgType, msgTypeToString(msgType), string(msg))
	}
}

func msgTypeToString(t int) string {
	switch t {
	case 1:
		return "TextMessage"
	case 2:
		return "BinaryMessage"
	case 8:
		return "CloseMessage"
	case 9:
		return "PingMessage"
	case 10:
		return "PongMessage"
	default:
		return "Unknown"
	}
}
