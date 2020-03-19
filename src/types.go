package main

type ConfigTypes struct {
	EmailSettings struct {
		Email    string `yaml:"email"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		To_email string `yaml:"to_email"`
	} `yaml:"email_settings"`
	MessageSettings struct {
		Sid      string `yaml:"sid"`
		Token    string `yaml:"token"`
		From_num string `yaml:"from_num"`
		To_num   string `yaml:"to_num"`
	} `yaml:"message_settings"`
}

type FailureMessage struct {
	Time         string `json:"time"`
	Failure_type string `json:"type"`
	Severity     int    `json:"severity"`
}

type MotionDetected struct {
	Time string
}

type IssueNotice struct {
	Severity  int `json:"severity"`
	Component string `json:"component"`
	Action    string `json:"action"`
}

type MonitorState struct {
	State bool
}

type RequestPower struct {
	Power     string `json:"power"`
	Severity  int    `json:"severity"`
	Component string `json:"component"`
}

type EventFH struct {
	Component    string
	Error_string string
	Time         string
	Severity     int
}

type MapMessage struct {
	message     string
	routing_key string
	time        string
	valid       bool
}

const FAILURE string = "Failure.*"
const FAILURENETWORK string = "Failure.Network"     //Level 4
const FAILUREDATABASE string = "Failure.Database"   //Level 4
const FAILURECOMPONENT string = "Failure.Component" //Level 2 if NAC 3
const FAILUREACCESS string = "Failure.Access"       //Level 5
const FAILURECAMERA string = "Failure.Camera"
const MOTIONDETECTED string = "Motion.Detected" //Level 5

const ISSUENOTICE string = "Issue.Notice"
const MONITORSTATE string = "Monitor.State"
const REQUESTPOWER string = "Request.Power"
const EVENTFH string = "Event.FH"
const EXCHANGENAME string = "topics"
const EXCHANGETYPE string = "topic"
const TIMEFORMAT string = "20060102150405"
const CAMERAMONITOR string = "CM"
const COMPONENT string = "FH"
const UPDATESTATEERROR string = "We have received a brand new state update"
const SERVERERROR string = "Server is failing to send"
const STATEUPDATESEVERITY int = 2
const SERVERSEVERITY int = 4
const FAILURECONVERT string = "Failed to convert"
const FAILUREPUBLISH string = "Failed to publish"

var SubscribedMessagesMap map[uint32]*MapMessage
var key_id uint32 = 0