package main

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

func messageFailure(issue bool) string {
	fail := ""
	if issue {
		fail = PublishEventFH(COMPONENT, SERVERERROR, getTime())
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
		log.Error("We have major flaw")
	}
	return shutdown_valid
}

func checkState() {
	for message_id := range SubscribedMessagesMap {
		if SubscribedMessagesMap[message_id].valid == true {
			log.Debug("Message id is: ", message_id)
			log.Debug("Message routing key is: ", SubscribedMessagesMap[message_id].routing_key)
			switch {
			case SubscribedMessagesMap[message_id].routing_key == MOTIONDETECTED:
				var message MotionDetected
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				messageFailure(SendEmailRoutine("We have movement in the flat", MOTIONMESSAGE))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILURENETWORK:
				log.Debug("Received a network failure message")
				messageFailure(SendEmailRoutine("Server unable to respond", "The network is not responding or the\n " +
												"firewall has shut down then network"))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILUREDATABASE:
				messageFailure(SendEmailRoutine("Data failure HouseGuard", "Serious Database failure"))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILURECOMPONENT:
				var message FailureMessage
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				log.Warn("Failure in component: ", message.Failure_type)
				messageFailure(SendEmailRoutine("Software not responding", "Serious Component failure, \n" +
					"please troubleshoot "  + message.Failure_type))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILUREACCESS:
				messageFailure(SendEmailRoutine("Multiple pin attempts", "Please check the alarm immediately"))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILURECAMERA:
				var message FailureMessage
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				valid := PublishRequestPower("restart", 5, CAMERAMONITOR)
				if valid != "" {
					log.Warn("Failed to publish")
				} else {
					log.Debug("Published Request Power")
					SubscribedMessagesMap[message_id].valid = false
				}

			case SubscribedMessagesMap[message_id].routing_key == MONITORSTATE:
				var monitor MonitorState
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &monitor)
				SetState(true)
				messageFailure(SendEmailRoutine(UPDATESTATE_TITLE, UPDATESTATE_MESSAGE))
				SetState(monitor.State)
				valid := PublishEventFH(COMPONENT, UPDATESTATE, getTime())
				if valid != "" {
					log.Warn("Failed to publish")
				} else {
					log.Debug("Published Event Fault Handler")
					SubscribedMessagesMap[message_id].valid = false
				}

			case SubscribedMessagesMap[message_id].routing_key == DEVICEFOUND:
				var device DeviceFound
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &device)
				if device.Status == BLOCKED {
					messageFailure(SendEmailRoutine(DEVICE_TITLE, 
					DEVICEBLOCKED_MESSAGE + device.Device_name))
				} else if device.Status == UNKNOWN {
					messageFailure(SendEmailRoutine(DEVICE_TITLE, 
						DEVICEUNKNOWN_MESSAGE + device.Device_name))
				}
				SubscribedMessagesMap[message_id].valid = false				

			default:
				log.Warn("We were not expecting this message unvalidating: ",
					SubscribedMessagesMap[message_id].routing_key)
				SubscribedMessagesMap[message_id].valid = false
			}
		}
	}

}
