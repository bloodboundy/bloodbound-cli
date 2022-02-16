package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	playerID   string
	authHeader string

	gameID    string
	gameIndex uint32
)

func setID(id string) {
	playerID = id
	authHeader = "Bearer " + playerID
}

func doAction(jso interface{}) error {
	u := url.URL{Scheme: "http", Host: *ADDR, Path: fmt.Sprintf("/games/%v/actions", gameID)}

	jsb, err := json.Marshal(jso)
	if err != nil {
		return errors.Wrap(err, "marshal req")
	}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsb))
	if err != nil {
		return errors.Wrap(err, "new req")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "post")
	}
	defer rsp.Body.Close()

	rspBytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return errors.Wrap(err, "read rsp")
	}
	logrus.Infof("request(%v): %v", string(rspBytes))
	return nil
}
