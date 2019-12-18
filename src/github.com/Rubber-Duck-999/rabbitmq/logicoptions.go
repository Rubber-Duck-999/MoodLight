package rabbitmq

import (
	"encoding/json"

	"github.com/Rubber-Duck-999/message"

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
	message.SetSettings(email, password, email, from_name, to_email)
	setup_invalid := message.TestEmail()
	log.Debug("Email test success : ", !setup_invalid)
	if setup_invalid == true {
		shutdown_valid = true
		messageFailure(shutdown_valid)
		log.Fatal("We have major flaw - shutting down node and diagonose")
	}
	return shutdown_valid
}

func checkState() {
	for message_id := range SubscribedMessagesMap {
		if SubscribedMessagesMap[message_id].valid == true {
			log.Debug("Message id is: ", message_id)
			log.Debug("Message routing key is: ", SubscribedMessagesMap[message_id].routing_key)
			switch {
			case SubscribedMessagesMap[message_id].routing_key == FAILURENETWORK:
				log.Debug("Received a network failure message")
				messageFailure(message.SendEmailRoutine("Serious Network failure"))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILUREDATABASE:
				messageFailure(message.SendEmailRoutine("Serious Database failure"))
				messageFailure(message.SendSMS("Serious Database failure"))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILURECOMPONENT:
				messageFailure(message.SendEmailRoutine("Serious Component failure"))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILUREACCESS:
				messageFailure(message.SendEmailRoutine("Serious Access Violation"))
				messageFailure(message.SendSMS("Serious Access Violation"))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILURECAMERA:
				var message FailureMessage
				log.Debug("Creating temp")
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				log.Debug("Converting json data")
				valid := PublishRequestPower("restart", message.Severity, CAMERAMONITOR)
				if valid != "" {
					SubscribedMessagesMap[message_id].valid = false
					log.Warn("Failed to publish")
				} else {
					log.Debug("Published Request Power")
				}

			case SubscribedMessagesMap[message_id].routing_key == ISSUENOTICE:
				var message IssueNotice
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				valid := PublishRequestPower(message.action, message.severity, message.component)
				if valid != "" {
					SubscribedMessagesMap[message_id].valid = false
					log.Warn("Failed to publish")
				} else {
					log.Debug("Published Request Power")
				}

			case SubscribedMessagesMap[message_id].routing_key == MONITORSTATE:
				var monitor MonitorState
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &monitor)
				message.SetState(monitor.state)
				valid := PublishEventFH(COMPONENT, UPDATESTATEERROR, getTime(), STATEUPDATESEVERITY)
				if valid != "" {
					SubscribedMessagesMap[message_id].valid = false
					log.Warn("Failed to publish")
				} else {
					log.Debug("Published Event Fault Handler")
				}

			default:
				log.Debug("We were not expecting this message unvalidating")
				SubscribedMessagesMap[message_id].valid = false
			}
		}
	}

}
