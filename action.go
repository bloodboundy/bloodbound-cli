package main

import (
	"strconv"

	"github.com/pkg/errors"
)

type ActionHandler func(t string, args ...string) error

var actionMap = map[string]ActionHandler{
	"target": targetAction,
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
