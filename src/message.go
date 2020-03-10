package main

import (
	"github.com/scorredoira/email"
	"github.com/sfreiberg/gotwilio"
	log "github.com/sirupsen/logrus"
	"net/mail"
	"net/smtp"
)

var _state bool
var _email string
var _password string
var _subject string
var _body string
var _from_email string
var _from_name string
var _to_email string
var _sid string
var _token string
var _from_num string
var _to_num string

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
	_token = ""
	_sid = ""
	_from_num = ""
	_to_num = ""
}

func SetState(state bool) {
	log.Debug("Requested to change our monitoring state")
	log.WithFields(log.Fields{
		"State": _state, "New State": state,
	}).Debug("State change")
	_state = state
}

func getState() bool {
	return _state
}

func SetSettings(email string, password string, from_email string,
	from_name string, to_email string) {
	log.Trace("Setting settings")
	_subject = "Test Email"
	_body = ""
	_email = email
	_password = password
	_from_email = from_email
	_from_name = from_name
	_to_email = to_email
}

func SetMessageSettings(sid string, token string, from_num string, to_num string) {
	_sid = sid
	_token = token
	_from_num = from_num
	_to_num = to_num
}

func TestEmail() bool {
	_subject = "Test Email"
	_body = ""
	fatal := sendEmail("Test")
	return fatal
}

func SendSMS(issue string) bool {
	state := false
	if _state {
		log.Debug("Sending important SMS")
		twilio := gotwilio.NewTwilioClient(_sid, _token)

		message := "Welcome to gotwilio!"
		_, _, err := twilio.SendSMS(_from_num, _to_num, message, "", "")
		if err != nil {
			return state
		} else {
			state = true
		}
	}
	return state
}

func SendEmailRoutine(issue string) bool {
	event := sendEmail(issue)
	return event
}

func sendEmail(issue string) bool {
	// compose the message
	fatal := false
	if _state {
		_body = issue
		m := email.NewMessage(_subject, _body)
		m.From = mail.Address{Name: _from_name, Address: _from_email}
		m.To = []string{_to_email}

		// send it
		auth := smtp.PlainAuth("", _email, _password, "smtp.zoho.eu")
		if err := email.Send("smtp.zoho.eu:587", auth, m); err != nil {
			log.Warn("Found a issue")
			log.Warn(err)
			fatal = true
		}
	}
	return fatal
}
