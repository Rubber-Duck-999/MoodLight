package main

import (
	"os"

	"github.com/akamensky/argparse"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Warn("FH - Beginning to run Fault Handler Program")
	parser := argparse.NewParser("file", "Config file for runtime purpose")
	// Create string flag
	f := parser.String("f", "config", &argparse.Options{Required: true, Help: "Necessary config"})
	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		log.Error(parser.Usage(err))
		os.Exit(2)
	}

	file := *f
	var data ConfigTypes
	if Exists(file) {
		GetData(&data, file)
	} else {
		log.Error("File doesn't exist")
		os.Exit(2)
	}
	log.Trace(data.EmailSettings.Email)
	if data.EmailSettings.Email != "" {
		SetEmailSettings(data.EmailSettings.Email,
			data.EmailSettings.Password,
			data.EmailSettings.Name,
			data.EmailSettings.To_email)
		SetMessageSettingsLogic(data.MessageSettings.Sid,
			data.MessageSettings.Token,
			data.MessageSettings.From_num,
			data.MessageSettings.To_num)
		if TestEmail() == true {
			log.Error("Cannot start test")
			os.Exit(1)
		}
		Subscribe()
	} else {
		log.Error("File not converted correctly")
		os.Exit(2)
	}
}
