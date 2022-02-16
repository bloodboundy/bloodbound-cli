package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ActionHandler func(t string, args ...string) error

var actionMap = map[string]ActionHandler{
	"target": targetAction,
	"new":    metaNewGame,
	"join":   metaJoinGame,
}

func execAction(args ...string) error {
	if len(args) < 1 {
		return errors.New("action type required")
	}
	t := args[0]
	h, ok := actionMap[t]
	if !ok {
		return errors.Errorf("unrecognized action type: %v", t)
	}
	return h(t, args[1:]...)
}

type actionJSONComm struct {
	Type     string `json:"type"`
	Operator string `json:"operator"`
	Round    uint32 `json:"round"`
	From     uint32 `json:"from"`
}

func makeActionJSONComm(t string) actionJSONComm {
	return actionJSONComm{
		Type:     t,
		Operator: playerID,
		From:     gameIndex,
	}
}

type TargetActionJSON struct {
	actionJSONComm
	To uint32 `json:"to"`
}

func targetAction(t string, args ...string) error {
	if len(args) < 1 {
		return errors.Errorf("args: to")
	}
	to, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "parse to")
	}
	return doAction(TargetActionJSON{makeActionJSONComm(t), uint32(to)})
}

type GameJSON struct {
	// meta data
	ID        string `json:"ID,omitempty"`
	CreatedAt uint64 `json:"created_at,omitempty"`
	Owner     string `json:"owner,omitempty"`

	// settings
	MaxPlayers *uint32 `json:"max_players,omitempty"`
	IsPrivate  bool    `json:"is_private,omitempty"`
	Password   string  `json:"password,omitempty"`
}

func metaNewGame(t string, args ...string) error {
	u := url.URL{Scheme: "http", Host: *ADDR, Path: "/games"}
	rsp, err := requestJSON("POST", u.String(), GameJSON{})
	if err != nil {
		return errors.Wrap(err, "new game")
	}
	logrus.Info("new game:", rsp)
	g := GameJSON{}
	if err := json.Unmarshal([]byte(rsp), &g); err != nil {
		return errors.Wrap(err, "unmarshal rsp")
	}
	gameID = g.ID
	return nil
}

func metaJoinGame(t string, args ...string) error {
	gid := gameID
	if len(args) > 0 {
		gid = args[0]
	}
	u := url.URL{Scheme: "http", Host: *ADDR, Path: fmt.Sprintf("/games/%v/players", gid)}
	rsp, err := requestJSON("POST", u.String(), struct {
		ID string `json:"id"`
	}{playerID})
	if err != nil {
		return errors.Wrap(err, "join game")
	}
	logrus.Info("join game:", rsp)
	return nil
}
