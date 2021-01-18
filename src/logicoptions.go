package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func SetEmailSettings(email string, password string, from_name string, to_email string) bool {
	shutdown_valid := false
	log.Trace("Email is: ", email)
	SetSettings(email, password, from_name, to_email)
	setup_invalid := TestEmail()
	log.Debug("Email test success : ", !setup_invalid)
	if setup_invalid {
		shutdown_valid = true
		log.Error("We have major flaw")
	}
	return shutdown_valid
}

func checkLogicMonitor(monitor MonitorState) string {
	var message string
	var state string
	if monitor.State {
		message = ACTIVATE_TITLE
		state = "ON"
		publishCameraStart()
	} else {
		state = "OFF"
		message = DEACTIVATE_TITLE
		publishCameraStop()
	}
	sendEmail(message, ACT_MESSAGE)
	valid := publishAlarmEvent("Admin", state)
	return valid
}

func cleanUp() {

	dirname := "." + string(filepath.Separator)

	d, err := os.Open(dirname)
	if err != nil {
		log.Warn(err)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		log.Warn(err)
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".jpg" {
				os.Remove(file.Name())
				log.Warn("Deleted ", file.Name())
			}
		}
	}
}

func checkState() {
	for message_id := range SubscribedMessagesMap {
		if SubscribedMessagesMap[message_id].valid {
			log.Debug("Message routing key is: ", SubscribedMessagesMap[message_id].routing_key)
			switch {
			case SubscribedMessagesMap[message_id].routing_key == MOTIONRESPONSE:
				log.Debug("Received a Motion Response Topic")
				var message MotionResponse
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				sendEmail("We have movement in the flat", MOTIONMESSAGE)
				cleanUp()
				publishMotion()
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILURENETWORK:
				log.Debug("Received a network failure message")
				sendEmail("Server unable to respond", "The network is not responding or the\n "+
					"firewall has shut down then network")
				status.LastFault = FAILURENETWORK
				StatusCheck()
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILURECOMPONENT:
				var message FailureMessage
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				log.Warn("Failure in component: ", message.Failure_type)
				sendEmail("Software not responding", "Serious Component failure, \n"+
					"please troubleshoot this issue: "+message.Failure_type)
				status.LastFault = FAILURECOMPONENT
				StatusCheck()
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILUREACCESS:
				sendEmail("Multiple pin attempts", "Please check the alarm immediately")
				status.LastFault = FAILUREACCESS
				StatusCheck()
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == FAILURECAMERA:
				status.LastFault = FAILURECAMERA
				StatusCheck()
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == MONITORSTATE:
				var monitor MonitorState
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &monitor)
				valid := checkLogicMonitor(monitor)
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
					sendEmail(DEVICE_TITLE,
						DEVICEBLOCKED_MESSAGE+device.Device_name+IPADDRESS+device.Ip_address)
				} else if device.Status == UNKNOWN {
					sendEmail(DEVICE_TITLE,
						DEVICEUNKNOWN_MESSAGE+device.Device_name+IPADDRESS+device.Ip_address)
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
