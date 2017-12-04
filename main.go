package main

import (
	"fmt"
	"os"

	"github.com/docker/go-plugins-helpers/sdk"
	"github.com/sirupsen/logrus"
	"gitlab.com/michael.golfi/appinsights/handler"
)

var logLevels = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
}

func main() {
	levelVal := os.Getenv("LOG_LEVEL")
	if levelVal == "" {
		levelVal = "info"
	}
	if level, exists := logLevels[levelVal]; exists {
		logrus.SetLevel(level)
	} else {
		fmt.Fprintln(os.Stderr, "invalid log level: ", levelVal)
		os.Exit(1)
	}

	h := sdk.NewHandler(`{"Implements": ["LoggingDriver"]}`)
	handler.Handle(&h, handler.NewDriver())
	if err := h.ServeUnix("appinsights", 0); err != nil {
		panic(err)
	}
}
