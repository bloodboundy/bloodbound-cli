package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type PlayerJSON struct {
	ID       string `json:"ID,omitempty"`
	Nickname string `json:"nickname,omitempty"`
}

func register() (id, nickname string, err error) {
	u := url.URL{Scheme: "http", Host: *ADDR, Path: "/register"}
	q := u.Query()
	q.Add("nickname", strconv.FormatUint(r.Uint64(), 16))
	u.RawQuery = q.Encode()
	rsp, err := http.Get(u.String())
	if err != nil {
		return "", "", errors.Wrap(err, "get")
	}
	d, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return "", "", errors.Wrap(err, "read")
	}
	p := &PlayerJSON{}
	if err := json.Unmarshal(d, p); err != nil {
		return "", "", errors.Wrap(err, "json.unmarshal")
	}
	return p.ID, p.Nickname, nil
}
