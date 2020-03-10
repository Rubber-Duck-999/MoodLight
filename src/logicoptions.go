package main

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

func messageFailure(issue bool) string {
	fail := ""
	if issue {
		fail = PublishEventFH(COMPONENT, SERVERERROR, getTime(), SERVERSEVERITY)
	}
	return fail
}

func SetEmailSettings(email string, password string, from_name string, to_email string) bool {
	shutdown_valid := false
	log.Trace("Email is: ", email)
	SetSettings(email, password, email, from_name, to_email)
	setup_invalid := TestEmail()
	log.Debug("Email test success : ", !setup_invalid)
	if setup_invalid {
		shutdown_valid = true
		messageFailure(shutdown_valid)
		log.Fatal("We have major flaw - shutting down node and diagonose")
	}
	return shutdown_valid
}

func SetMessageSettingsLogic(sid string, token string, from_num string, to_num string) {
	SetMessageSettings(sid, token, from_num, to_num)
}

func checkState() {
	for message_id := range SubscribedMessagesMap {
		if SubscribedMessagesMap[message_id].valid == true {
			log.Debug("Message id is: ", message_id)
			log.Debug("Message routing key is: ", SubscribedMessagesMap[message_id].routing_key)
			switch {
			case SubscribedMessagesMap[message_id].routing_key == FAILURENETWORK:
				log.Debug("Received a network failure message")
				messageFailure(SendEmailRoutine("Serious Network failure"))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILUREDATABASE:
				messageFailure(SendEmailRoutine("Serious Database failure"))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILURECOMPONENT:
				messageFailure(SendEmailRoutine("Serious Component failure"))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILUREACCESS:
				messageFailure(SendEmailRoutine("Serious Access Violation"))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILURECAMERA:
				var message FailureMessage
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				valid := PublishRequestPower("restart", message.Severity, CAMERAMONITOR)
				if valid != "" {
					SubscribedMessagesMap[message_id].valid = false
					log.Warn("Failed to publish")
				} else {
					log.Debug("Published Request Power")
				}

			case SubscribedMessagesMap[message_id].routing_key == ISSUENOTICE:
				var message IssueNotice
				err := json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				if err == nil {
					valid := PublishRequestPower("restart", message.Severity, message.Component)
					log.Info("We will inform them to shutdown: ", message.Component)
					if valid != "" {
						log.Warn("Failed to publish")
					} else {
						SubscribedMessagesMap[message_id].valid = false
						log.Info("Published Request Power")
					}
				} else {
					log.Warn(err)
				}

			case SubscribedMessagesMap[message_id].routing_key == MONITORSTATE:
				var monitor MonitorState
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &monitor)
				SetState(monitor.State)
				valid := PublishEventFH(COMPONENT, UPDATESTATEERROR, getTime(), STATEUPDATESEVERITY)
				if valid != "" {
					log.Warn("Failed to publish")
				} else {
					log.Debug("Published Event Fault Handler")
					SubscribedMessagesMap[message_id].valid = false
				}

			case SubscribedMessagesMap[message_id].routing_key == MOTIONDETECTED:
				messageFailure(SendEmailRoutine("Motion Detected"))
				messageFailure(SendSMS("Motion Detected"))
				SubscribedMessagesMap[message_id].valid = false

			default:
				log.Warn("We were not expecting this message unvalidating: ",
					SubscribedMessagesMap[message_id].routing_key)
				SubscribedMessagesMap[message_id].valid = false
			}
		}
	}

}
