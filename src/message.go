package main

import (
	"time"

	"net/mail"
	"net/smtp"

	"github.com/scorredoira/email"
	log "github.com/sirupsen/logrus"
)

var _state bool
var _email string
var _password string
var _subject string
var _body string
var _from_email string
var _from_name string
var _to_email string
var _year int
var _month time.Month
var _day int
var _messages_sent int

func init() {
	log.Trace("Initialised message package")
	_state = false
	_subject = ""
	_body = ""
	_email = ""
	_password = ""
	_from_email = ""
	_from_name = ""
	_to_email = ""
	setDate()
}

func SetState(state bool) {
	log.Debug("Requested to change our monitoring state")
	log.Debug("State change from: ", _state, " to: ", state)
	_state = state
}

func getState() bool {
	return _state
}

func SetSettings(email string, password string, from_email string,
	from_name string) {
	log.Trace("Setting settings")
	_subject = "Test Email"
	_body = ""
	_email = email
	_password = password
	_from_email = from_email
	_from_name = from_name
	_to_email = from_email
}

func TestEmail() bool {
	_subject = "Test Email"
	_body = ""

	fatal := sendEmail("Starting up Server", "Test")
	sendLogsEmail()
	return fatal
}

func sendLogsEmail() {
	log.Debug("Sending logs email")
	m := email.NewMessage("HouseGuard Daily Logs", "All included")
	m.From = mail.Address{Name: _from_name, Address: _from_email}
	m.To = []string{_to_email}

	//Attachments
	var files = [7]string{"logs/DBM.txt", "logs/EVM.txt", "logs/oldFH.txt",
		"logs/NAC.txt", "logs/SYP.txt", "logs/UP.txt",
		"logs/cameramonitor.log"}
	for _, file := range files {
		if Exists(file) {
			if err := m.Attach(file); err != nil {
				log.Error(err)
			}
		}
	}

	// send it
	auth := smtp.PlainAuth("", _email, _password, "smtp.zoho.eu")
	if emailErr := email.Send("smtp.zoho.eu:587", auth, m); emailErr != nil {
		log.Warn("Found a issue")
		log.Warn(emailErr)
	}
}

func setDate() {
	_year, _month, _day = time.Now().Date()
	_messages_sent = 0
}

func checkCanSend() bool {
	year, month, day := time.Now().Date()
	if year == _year {
		if month == _month {
			if day == _day {
				if _messages_sent <= 30 {
					_messages_sent++
					return true
				} else {
					log.Error("Max messages sent")
					return false
				}
			} else {
				setDate()
				checkCanSend()
			}
		}
	}
	return false
}

func sendEmail(subject string, issue string) bool {
	// compose the message
	fatal := false
	if _state && checkCanSend() {
		log.Debug("Sending email")
		_body = issue
		m := email.NewMessage(subject, _body)
		m.From = mail.Address{Name: _from_name, Address: _from_email}
		m.To = []string{_to_email}

		// send it
		auth := smtp.PlainAuth("", _email, _password, "smtp.zoho.eu")
		if err := email.Send("smtp.zoho.eu:587", auth, m); err != nil {
			log.Warn("Found a issue")
			log.Warn(err)
			fatal = true
		}
	} else {
		log.Warn("We cannot send an email currently as state: ", _state)
	}
	return fatal
}

func sendAttachmentEmail(issue string, file string) bool {
	// compose the message
	fatal := false
	if _state && checkCanSend() {
		log.Debug("Sending email")
		_body = issue
		_subject = "Movement in Flat"
		m := email.NewMessage(_subject, _body)
		m.From = mail.Address{Name: _from_name, Address: _from_email}
		m.To = []string{_to_email}

		//Attachments
		if Exists(file) {
			if err := m.Attach(file); err != nil {
				log.Fatal(err)
				fatal = true
			}
		}

		// send it
		auth := smtp.PlainAuth("", _email, _password, "smtp.zoho.eu")
		if emailErr := email.Send("smtp.zoho.eu:587", auth, m); emailErr != nil {
			log.Warn("Found a issue")
			log.Warn(emailErr)
			fatal = true
		}
	}
	return fatal
}
