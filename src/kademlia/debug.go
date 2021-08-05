package kademlia

import (
	logrus "github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func LogInit()  {
	log, _ := os.Create("log.txt")
	Log.SetOutput(log)
	Log.SetReportCaller(true)
	Log.Formatter = &logrus.TextFormatter{
		ForceColors: true,
	}
}

/*
	Error
	Trace
	Info
	FATA
	WARN

	template:
	Log.WithFields(logrus.Fields{
		"xxx" : x,
	}).Error("233")
*/