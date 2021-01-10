package main

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

func messageFailure(issue bool) string {
	fail := ""
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

func checkLogicMonitor(monitor MonitorState) string {
	var message string
	var state string
	if monitor.State == true {
		message = ACTIVATE_TITLE
		state = "ON"
		publishCameraStart()
	} else {
		state = "OFF"
		message = DEACTIVATE_TITLE
		publishCameraStop()
	}
	messageFailure(sendEmail(message, ACT_MESSAGE))
	valid := publishAlarmEvent("Admin", state)
	return valid
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
			log.Debug("Message routing key is: ", SubscribedMessagesMap[message_id].routing_key)
			switch {
			case SubscribedMessagesMap[message_id].routing_key == MOTIONDETECTED:
				var message MotionDetected
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				if email_changed {
					messageFailure(sendEmail("We have movement in the flat", MOTIONMESSAGE))
					SubscribedMessagesMap[message_id].valid = false
				}

			case SubscribedMessagesMap[message_id].routing_key == FAILURENETWORK:
				log.Debug("Received a network failure message")
				if email_changed {
					messageFailure(sendEmail("Server unable to respond", "The network is not responding or the\n "+
						"firewall has shut down then network"))
					status.DailyFaults = checkDay(status.DailyFaults)
					network.Count++
					SubscribedMessagesMap[message_id].valid = false
				}

			case SubscribedMessagesMap[message_id].routing_key == FAILUREDATABASE:
				log.Debug("Received a database failure message")
				if email_changed {
					messageFailure(sendEmail("Data failure HouseGuard", "Serious Database failure"))
					status.DailyFaults = checkDay(status.DailyFaults)
					database.Count++
					SubscribedMessagesMap[message_id].valid = false
				}

			case SubscribedMessagesMap[message_id].routing_key == FAILURECOMPONENT:
				var message FailureMessage
				if email_changed {
					json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
					log.Warn("Failure in component: ", message.Failure_type)
					messageFailure(sendEmail("Software not responding", "Serious Component failure, \n"+
						"please troubleshoot this issue: "+message.Failure_type))
					status.DailyFaults = checkDay(status.DailyFaults)
					software.Count++
					SubscribedMessagesMap[message_id].valid = false
				}

			case SubscribedMessagesMap[message_id].routing_key == FAILUREACCESS:
				if email_changed {
					messageFailure(sendEmail("Multiple pin attempts", "Please check the alarm immediately"))
					status.DailyFaults = checkDay(status.DailyFaults)
					access.Count++
					SubscribedMessagesMap[message_id].valid = false
				}

			case SubscribedMessagesMap[message_id].routing_key == FAILURECAMERA:
				camera.Count++
				status.DailyFaults = checkDay(status.DailyFaults)
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == MONITORSTATE:
				var monitor MonitorState
				if email_changed {
					json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &monitor)
					valid := checkLogicMonitor(monitor)
					if valid != "" {
						log.Warn("Failed to publish")
					} else {
						log.Debug("Published Event Fault Handler")
						SubscribedMessagesMap[message_id].valid = false
					}
				}

			case SubscribedMessagesMap[message_id].routing_key == DEVICEFOUND:
				var device DeviceFound
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &device)
				if email_changed {
					if device.Status == BLOCKED {
						messageFailure(sendEmail(DEVICE_TITLE,
							DEVICEBLOCKED_MESSAGE+device.Device_name+IPADDRESS+device.Ip_address))
					} else if device.Status == UNKNOWN {
						messageFailure(sendEmail(DEVICE_TITLE,
							DEVICEUNKNOWN_MESSAGE+device.Device_name+IPADDRESS+device.Ip_address))
					}
					SubscribedMessagesMap[message_id].valid = false
				}

			default:
				log.Warn("We were not expecting this message unvalidating: ",
					SubscribedMessagesMap[message_id].routing_key)
				SubscribedMessagesMap[message_id].valid = false
			}
			StatusCheck()
		}
	}

}
