package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func setupLogger() {
	var level logrus.Level
	err := level.UnmarshalText([]byte(*LEVEL))
	if err != nil {
		logrus.Panicf("unaccptable level: %v, accepts: trace debug info warning error fatal panic", *LEVEL)
	}
	logrus.SetLevel(level)

	out := os.Stderr
	if *OUTPUT != "" {
		file, err := os.OpenFile(*OUTPUT, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.Panicf("failed to open log output file: %v", err)
		}
		out = file
	}
	logrus.SetOutput(out)
}
