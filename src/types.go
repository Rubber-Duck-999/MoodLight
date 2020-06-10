package main

type ConfigTypes struct {
	EmailSettings struct {
		Email    string `yaml:"email"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		To_email string `yaml:"to_email"`
	} `yaml:"email_settings"`
	MessageSettings struct {
		Password string `yaml:"passcode"`
	} `yaml:"message_settings"`
}

type FailureMessage struct {
	Time         string `json:"time"`
	Failure_type string `json:"type_of_failure"`
}

type MotionDetected struct {
	File string `json:"file"`
	Time string `json:"time"`
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
	Component    string `json:"component"`
	Message      string `json:"message"`
	Time         string `json:"time"`
	EventTypeId  string `json:"event_type_id"`
}

type MapMessage struct {
	message     string
	routing_key string
	time        string
	valid       bool
}

type DeviceFound struct {
	Device_name string `json:"name"`
	Ip_address string `json:"address"`
	Status int `json:"status"`
}

type StatusFH struct {
	DailyFaults int `json:"daily_faults"`
	CommonFaults string `json:"common_faults"`
}

type Fault struct {
    Count int
    Name  string
}

type GUIDUpdate struct {
	GUID string `json:"guid"`
}

const FAILURE string = "Failure.*"
const FAILURENETWORK string = "Failure.Network"     //Level 5
const FAILUREDATABASE string = "Failure.Database"   //Level 4
const FAILURECOMPONENT string = "Failure.Component" //Level 3
const FAILUREACCESS string = "Failure.Access"       //Level 6
const FAILURECAMERA string = "Failure.Camera" // Level 2
const MOTIONDETECTED string = "Motion.Detected" //Level 7

const DEVICEFOUND string = "Device.Found"
const MONITORSTATE string = "Monitor.State"
const REQUESTPOWER string = "Request.Power"
const EVENTFH string = "Event.FH"
const STATUSFH string = "Status.FH"
const GUIDUPDATE string = "GUID.Update"
const EXCHANGENAME string = "topics"
const EXCHANGETYPE string = "topic"
const TIMEFORMAT string = "2006/01/02 15:04:05"
const CAMERAMONITOR string = "CM"
const COMPONENT string = "FH"
//
const DEVICE_TITLE string = "NETWORK ERROR"
const DEVICEBLOCKED_MESSAGE string = "A blocked dveice has joined the\n" +
									 "network. Device name: "
const DEVICEUNKNOWN_MESSAGE string = "A unknown device has joined the\n" +
									 "network. Device name: "
//
const UPDATESTATE_TITLE string = "Monitoring state changed"
const UPDATESTATE_MESSAGE string = "The alarm state has been changed \n" +
									"Please ensure that whoever enacted this " +
									"was authorised to do so"
//
const UPDATESTATE string = "Motion state changed"
const SERVERERROR string = "Server is failing to send"
const MOTIONMESSAGE string = "There was movement, check the image on the drive"
//
const GUIDUPDATE_TITLE string = "Daily GUID Key Inside"
const GUIDUPDATE_MESSAGE string = "Key: " 
//
const STATEUPDATESEVERITY int = 2
const SERVERSEVERITY int = 4
const BLOCKED int = 2
const UNKNOWN int = 4
const FAILURECONVERT string = "Failed to convert"
const FAILUREPUBLISH string = "Failed to publish"

var SubscribedMessagesMap map[uint32]*MapMessage
var key_id uint32 = 0
