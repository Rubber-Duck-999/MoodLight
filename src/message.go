package main

import (
	"time"

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
	setDate()
}

func TestEmail() bool {
	_subject = "Test Email"
	_body = ""
	fatal := sendEmail("Starting up Server", "Test")
	return fatal
}

func SendSMS(issue string) bool {
	state := false
	if _state && checkCanSend() {
		log.Debug("Sending important SMS")
		twilio := gotwilio.NewTwilioClient(_sid, _token)

		_, exception, err := twilio.SendSMS(_from_num, _to_num, issue, "","")
		if err != nil || exception != nil{
			log.Error("Exception found: ", exception)
			log.Error("SMS Failure: ", err)
			state = true
			return state
		}
	}
	return state
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

func SendEmailRoutine(subject string, issue string) bool {
	event := sendEmail(subject, issue)
	return event
}

func SendAttachedRoutine(issue string, file string) bool {
	event := sendAttachmentEmail(issue, file)
	return event
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
