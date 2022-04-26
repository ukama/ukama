package lwm2m

import (
	"errors"
	"lwm2m-gateway/pkg/senml"
	stat "lwm2m-gateway/specs/common/spec"
	spec "lwm2m-gateway/specs/lwm2mIface/spec"
	"strings"

	"strconv"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const (
	ALARM_OBJ_ID string = "34570"
)

var EvtCB func(proto.Message, string)

type AlarmObj struct {
	EType             int
	RealTime          bool
	State             uint32
	Disc              string
	LowThreshold      float32
	HighThreshold     float32
	CriticalThreshold float32
	EventCount        int
	Time              int
	ObjId             uint32
	InstId            uint32
	ResourceId        uint32
	SensorValue       float32
	SensorUnits       string
	ApplicationType   string
}

//Register Callback for Events.
func RegisterCallbackForEvents(cb func(proto.Message, string)) {
	EvtCB = cb
}

//Creare a alram Event
func CreateAlarmEvt(token uint32, uuid string, url string, alarm *AlarmObj) *spec.Lwm2MAlarmMsg {
	var limit float32

	// Add the limit value
	if alarm.SensorValue < alarm.LowThreshold {
		limit = alarm.LowThreshold
	} else if alarm.SensorValue > alarm.CriticalThreshold {
		limit = alarm.CriticalThreshold
	} else if alarm.SensorValue > alarm.HighThreshold && alarm.SensorValue < alarm.CriticalThreshold {
		limit = alarm.HighThreshold
	}

	// Create event Message for controller
	evtMsg := spec.Lwm2MAlarmMsg{
		Token:       uint64(token),
		Uuid:        uuid,
		EvtType:     int32(alarm.EType),
		State:       uint32(alarm.State),
		Alarmurl:    url,
		Resrcurl:    "/" + strconv.Itoa(int(alarm.ObjId)) + "/" + strconv.Itoa(int(alarm.InstId)) + "/" + strconv.Itoa(int(alarm.ResourceId)),
		SensorValue: alarm.SensorValue,
		SensorLimit: limit,
		SensorUnit:  alarm.SensorUnits,
		AckBy:       "",
	}

	log.Infof("Event:: Prepared Event %d for Device %s AlarmURL %s Resource URL %s for reporting to controller.", evtMsg.Token, evtMsg.Uuid, evtMsg.Alarmurl, evtMsg.Resrcurl)
	return &evtMsg
}

//Get SENML Json Value
func GetSENMLJsonAlarmValue(dType string, rec senml.Record, val interface{}) error {

	switch data := val.(type) {
	case *string:
		*data = *rec.StringValue
	case *int:
		*data = int(*rec.Value)
	case *uint32:
		*data = uint32(*rec.Value)
	case *float32:
		*data = float32(*rec.Value)
	case *bool:
		*data = *rec.BoolValue
	default:
		return errors.New("unknown data type in event message")
	}

	return nil

	//Read the value
	// switch dType {
	// case "String":
	// 	if rec.StringValue != nil {
	// 		org := val.(*string)
	// 		*org = *rec.StringValue
	// 		val = org
	// 		log.Debugf("Event:: SENML Value read is %s", *rec.StringValue)
	// 	}
	// case "Integer":
	// 	if rec.Value != nil {
	// 		org := val.(*int)
	// 		*org = int(*rec.Value)
	// 		val = org
	// 		log.Debugf("Event:: SENML Value read is %d", *rec.Value)
	// 	}
	// case "Float":
	// 	if rec.Value != nil {
	// 		val = *rec.Value
	// 		log.Debugf("Event:: SENML Value read is %f", *rec.Value)
	// 	}
	// case "Boolean":
	// 	if rec.BoolValue != nil {
	// 		val = *rec.BoolValue
	// 		log.Debugf("Event:: SENML Value read is %t", *rec.BoolValue)
	// 	}
	// default:
	// 	//TODO:: Check what do here
	// 	return errors.New("Unkown Data Type for ResourceID")
	// }
	// log.Debugf("Event:: SENML Interface Value read is %+v", val)
	// return nil

}

// Alarm Object
func AlarmObject(cfgJson *senml.Msg) (*AlarmObj, error) {

	//Get resource ID to data type map.
	rIDToTypeMap, err := GetResourceIDDataType(ALARM_OBJ_ID)
	if err != nil {
		log.Errorf("Event:: Error getting ResourceID to data type map:: Error %s", err)
		return nil, err
	}

	alarm := AlarmObj{}

	//Update Config data now
	records := cfgJson.Records
	for _, rec := range records {
		if rec.Name != nil {
			//convert resourceID to integer
			rid, err := strconv.Atoi(*rec.Name)
			if err != nil {
				return nil, errors.New("Couldn't decode ResourceID from senml record.")
			}

			dType := rIDToTypeMap[ID(rid)]

			switch rid {
			case 6011:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.EType)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 6011, err.Error(), rec)
					return nil, err
				}
			case 6012:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.RealTime)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 6012, err.Error(), rec)
					return nil, err
				}
			case 1:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.State)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 1, err.Error(), rec)
					return nil, err
				}
			case 5:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.LowThreshold)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 5, err.Error(), rec)
					return nil, err
				}
			case 6:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.HighThreshold)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 6, err.Error(), rec)
					return nil, err
				}
			case 7:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.CriticalThreshold)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 6, err.Error(), rec)
					return nil, err
				}
			case 8:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.ObjId)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 8, err.Error(), rec)
					return nil, err
				}
			case 9:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.InstId)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 9, err.Error(), rec)
					return nil, err
				}
			case 10:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.ResourceId)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 10, err.Error(), rec)
					return nil, err
				}
			case 13:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.Disc)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 13, err.Error(), rec)
					return nil, err
				}
			case 6018:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.EventCount)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 6018, err.Error(), rec)
					return nil, err
				}
			case 6021:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.Time)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 6021, err.Error(), rec)
					return nil, err
				}
			case 5700:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.SensorValue)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 5700, err.Error(), rec)
					return nil, err
				}
			case 5701:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.SensorUnits)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 5701, err.Error(), rec)
					return nil, err
				}
			case 5750:
				err := GetSENMLJsonAlarmValue(dType, rec, &alarm.ApplicationType)
				if err != nil {
					log.Errorf("Event:: Error getting GetSENMLJsonAlarmValue %d Error %s rec %+v", 5750, err.Error(), rec)
					return nil, err
				}
			}
		}
	}
	log.Debugf("Event:: Alarm object we read is %+v.", alarm)
	return &alarm, nil
}

// Process Alarm Data
func processReadAlarmData(resp *ResponseMsg) (*AlarmObj, error) {

	// Decode ContentType
	cfgJson, err := decodeContent(Lwm2mContentType(resp.Format), resp.Message)
	if err != nil {
		log.Errorf("Event:: Problem while decoding SENML-JSON Reply from LWM2M server.")
		return nil, err
	}
	log.Debugf("Event:: Alarm Data for %d decoded.", resp.Token)
	return AlarmObject(cfgJson)

}

// Adjust URL to object/instance
func GetAlarmObjectURL(ourl string) *string {
	var newUrl string
	s := strings.Split(ourl, "/")
	if len(s) > 2 {
		newUrl = "/" + s[1] + "/" + s[2]
	}
	log.Debugf("Event:: URL is %s %+v", newUrl, s)

	return &newUrl
}

// Reading event details
func readAlarm(token uint32, uuid string, nurl string) *spec.Lwm2MAlarmMsg {
	log.Infof("Event:: Preparing Alarm data read request [%d] for device %s URL %s (length %d) ", token, uuid, nurl, len(nurl))

	url := GetAlarmObjectURL(nurl)
	if url == nil {
		log.Errorf("Event:: Unable to parse URL from the event %d for device %s. URL was [%s].", token, uuid, nurl)
		return nil
	}

	//Prepare message
	msg := PrepareMsgForLwm2mServer(Read, &uuid, url, nil)
	log.Infof("Event:: Requesting Alarm data from Lwm2m Server %d and Data is %s", msg.Token, msg.Message)

	//Send and wait for response.
	resp, err := transmit(uuid, *msg)
	if err == nil {

		// handle response
		if resp.Status == COAP_205_CONTENT {

			// Decode the responded config data from LwM2M server.
			alarm, err := processReadAlarmData(resp)
			if (err == nil) && (alarm != nil) {
				alarmEvt := CreateAlarmEvt(token, uuid, *url, alarm)
				log.Debugf("Event:: Event:: Reporting event to controller %+v.", alarmEvt)
				return alarmEvt
			} else {
				log.Errorf("Event:: Failed to process read response from  LwM2M server. Error :: %+v", err)
			}

		}
		log.Errorf("Event:: Event:: Invalid response %s from LwM2M server for %s from device %s", stat.StatusCode_name[int32(resp.Status)], *url, uuid)
	} else {
		log.Errorf("Event:: Event:: Failed to request LwM2M server for %s from device %s. Err %s", *url, uuid, err.Error())
	}
	return nil

}

//Hanlding event
func handleEvent(evt *EventMsg) {
	log.Debugf("Event:: Event:: Handling Event %d for %s URL %s ", evt.Token, evt.Uuid, evt.Uri)
	EvtAlarm := readAlarm(evt.Token, evt.Uuid, evt.Uri)
	if EvtAlarm != nil {
		EvtCB(EvtAlarm, evt.Uuid)
	}
}
