package main

import (
	"bufio"
	"flag"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/google/shlex"
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
	setID(id)

	u := url.URL{Scheme: "ws", Host: *ADDR, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logrus.Panicf("ws.dial: %v", err)
	}
	defer c.Close()

	output := make(chan []byte)
	done := make(chan struct{})
	go readFromWS(c, output, done)

	handleStdIn()
}

func handleStdIn() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		in := scanner.Text()
		logrus.Infof("input: %v", in)

		tokens, err := shlex.Split(in)
		if err != nil {
			logrus.Errorf("shlex: %v", err)
			continue
		}
		if err = execAction(tokens...); err != nil {
			logrus.Error(err)
		}
	}
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
