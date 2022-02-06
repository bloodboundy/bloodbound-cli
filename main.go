package main

import (
	"flag"
	"math/rand"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var ADDR = flag.String("addr", "localhost:8080", "server location")

// Log flags
var (
	LEVEL = flag.String("log-level", "info",
		`log level, accepts: trace debug info warning error fatal panic, default: info`)
	OUTPUT = flag.String("log-output", "", `the file log will be output to, default to stderr`)
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	flag.Parse()
	setupLogger()

	id, nickname, err := register()
	if err != nil {
		logrus.Panicf("register: %v", err)
	}
	logrus.Infof("nickname: %v, ID: %v", nickname, id)

	u := url.URL{Scheme: "ws", Host: *ADDR, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logrus.Panicf("ws.dial: %v", err)
	}
	defer c.Close()

	output := make(chan []byte)
	done := make(chan struct{})
	go readFromWS(c, output, done)
}

func readFromWS(c *websocket.Conn, output chan []byte, done chan struct{}) {
	defer close(done)
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			logrus.Errorf("ws.read: %v", err)
			return
		}
		logrus.Debugf("ws.read: %v", string(msg))
		output <- msg
	}
}
