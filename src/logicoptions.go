package main

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

func messageFailure(issue bool) string {
	fail := ""
	if issue {
		fail = PublishEventFH(COMPONENT, SERVERERROR, getTime(), "FH1")
	}
	return fail
}

func checkDay(daily int) int {
	_, _, current_day := time.Now().Date()
	if current_day == day {
		daily++
		return daily
	} else {
		_, _, day = time.Now().Date()
		return 0
	}
}

func SetEmailSettings(email string, password string, from_name string) bool {
	shutdown_valid := false
	log.Trace("Email is: ", email)
	SetSettings(email, password, email, from_name)
	setup_invalid := TestEmail()
	log.Debug("Email test success : ", !setup_invalid)
	if setup_invalid {
		shutdown_valid = true
		messageFailure(shutdown_valid)
		log.Error("We have major flaw")
	}
	return shutdown_valid
}

func GetCommonFault() (string, int) {
	max := 0
	fault_string := "None"
	faults := []Fault{network, database, software,
		access, camera}
	for _, local := range faults {
		if local.Count > max {
			max = local.Count
			fault_string = local.Name
		}
	}
	log.Debug("Common Fault found: ", fault_string+
		" with ", max, " faults")
	return fault_string, max
}

func checkState() {
	for message_id := range SubscribedMessagesMap {
		if SubscribedMessagesMap[message_id].valid == true {
			if first {
				PublishEmailRequest(ADMIN_ROLE)
				first = false
			}
			log.Debug("Message id is: ", message_id)
			log.Debug("Message routing key is: ", SubscribedMessagesMap[message_id].routing_key)
			switch {
			case SubscribedMessagesMap[message_id].routing_key == MOTIONDETECTED:
				var message MotionDetected
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				messageFailure(sendEmail("We have movement in the flat", MOTIONMESSAGE, ""))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == EMAILRESPONSE:
				var message EmailResponse
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				for _, account := range message.Accounts {
					log.Debug("Received: ", account.role, " and email: ", account.email)
					if account.role == ADMIN_ROLE {
						_to_email = account.email
					}
				}
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILURENETWORK:
				log.Debug("Received a network failure message")
				messageFailure(sendEmail("Server unable to respond", "The network is not responding or the\n "+
					"firewall has shut down then network", ""))
				status.DailyFaults = checkDay(status.DailyFaults)
				network.Count++
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILUREDATABASE:
				log.Debug("Received a database failure message")
				messageFailure(sendEmail("Data failure HouseGuard", "Serious Database failure", ""))
				status.DailyFaults = checkDay(status.DailyFaults)
				database.Count++
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILURECOMPONENT:
				var message FailureMessage
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				log.Warn("Failure in component: ", message.Failure_type)
				messageFailure(sendEmail("Software not responding", "Serious Component failure, \n"+
					"please troubleshoot this issue: "+message.Failure_type, ""))
				status.DailyFaults = checkDay(status.DailyFaults)
				software.Count++
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILUREACCESS:
				messageFailure(sendEmail("Multiple pin attempts", "Please check the alarm immediately", ""))
				status.DailyFaults = checkDay(status.DailyFaults)
				access.Count++
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILURECAMERA:
				camera.Count++
				status.DailyFaults = checkDay(status.DailyFaults)
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == GUIDUPDATE:
				var guidUpdate GUIDUpdate
				log.Debug("GUID Update")
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &guidUpdate)
				messageFailure(sendEmail(GUIDUPDATE_TITLE, GUIDUPDATE_MESSAGE+guidUpdate.GUID, ""))
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == MONITORSTATE:
				var monitor MonitorState
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &monitor)
				SetState(true)
				messageFailure(sendEmail(UPDATESTATE_TITLE, UPDATESTATE_MESSAGE, ""))
				SetState(monitor.State)
				valid := PublishEventFH(COMPONENT, UPDATESTATE, getTime(), "FH2")
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
					messageFailure(sendEmail(DEVICE_TITLE,
						DEVICEBLOCKED_MESSAGE+device.Device_name, ""))
				} else if device.Status == UNKNOWN {
					messageFailure(sendEmail(DEVICE_TITLE,
						DEVICEUNKNOWN_MESSAGE+device.Device_name, ""))
				}
				SubscribedMessagesMap[message_id].valid = false

			default:
				log.Warn("We were not expecting this message unvalidating: ",
					SubscribedMessagesMap[message_id].routing_key)
				SubscribedMessagesMap[message_id].valid = false
			}
			StatusCheck()
		}
	}

}
