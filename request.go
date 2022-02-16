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
	rsp, err := requestJSON("POST", u.String(), bytes.NewBuffer(jsb))
	if err != nil {
		return errors.Wrap(err, "new req")
	}

	logrus.Infof("action(%s): %v", jsb, string(rsp))
	return nil
}

func requestJSON(method string, url string, jso interface{}) (response string, err error) {
	jsb, err := json.Marshal(jso)
	if err != nil {
		return "", errors.Wrap(err, "marshal req")
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsb))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	rspBytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read rsp")
	}
	return string(rspBytes), nil
}
