package lwm2m

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	cfg "lwm2m-gateway/pkg/config"
	"lwm2m-gateway/pkg/senml"
	stat "lwm2m-gateway/specs/common/spec"
	spec "lwm2m-gateway/specs/lwm2mIface/spec"
	"reflect"
	"strconv"

	toml "github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
	"github.com/theherk/viper"
)

// Lwm2m Srever expects request contetnt to be string.
type LwM2MServerReq struct {
	//msg string
}

//LwM2M server response could be in multiple content type
type LwM2MServerResp struct {
	//msg         string
	//uri         spec.URI
	//contentType Lwm2mContentType
	//length      int
	//data        string // TODO::Change to []byte
	//status      int32
}

type ResourceIDTypeMap struct {
	//resourceid ID
	//dType      string
}

// Generic function to check if item exist in any array
func itemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

// Validate if object id
func IsValidObjectId(objectid ID) bool {
	return itemExists(ObjectIDList, objectid)
}

// URI to string
func URItoString(uri *spec.URI, op Lwm2mop) *string {
	if IsValidObjectId(ID(uri.Object)) {
		uriext := "/" + strconv.Itoa(int(uri.Object)) + "/" + strconv.Itoa(int(uri.Instance))
		//This is because proto3 by default set int value to 3 even if variable is not transmitted in wire.
		if uri.Resource == 0 {
			uriext = uriext + ""
		} else {
			uriext = uriext + "/" + strconv.Itoa(int(uri.Resource))
		}
		log.Debugf("LwM2M::URI formed for %d %d %d is %s", uri.Object, uri.Instance, uri.Resource, uriext)

		return &uriext
	} else {
		return nil
	}
}

//Translate response codes
func LwM2MResponseCode(code uint32) stat.StatusCode {

	var status stat.StatusCode
	switch code {
	case COAP_NO_ERROR:
		status = stat.StatusCode_STATUS_OK
	case COAP_201_CREATED:
		status = stat.StatusCode_STATUS_CREATED
	case COAP_202_DELETED:
		status = stat.StatusCode_STATUS_DELETED
	case COAP_204_CHANGED:
		status = stat.StatusCode_STATUS_CHANGED
	case COAP_205_CONTENT:
		status = stat.StatusCode_STATUS_CONTENT
	case COAP_400_BAD_REQUEST:
		status = stat.StatusCode_ERR_BADREQUEST
	case COAP_401_UNAUTHORIZED:
		status = stat.StatusCode_ERR_UNAUTORISED
	case COAP_404_NOT_FOUND:
		status = stat.StatusCode_ERR_NOT_FOUND
	case COAP_405_METHOD_NOT_ALLOWED:
		status = stat.StatusCode_ERR_METHOD_NOT_ALLOWED
	case COAP_406_NOT_ACCEPTABLE:
		status = stat.StatusCode_ERR_METHOD_NOT_ACCEPTABLE
	case COAP_500_INTERNAL_SERVER_ERROR:
		status = stat.StatusCode_ERR_INTERNAL_SERVER_ERROR
	case COAP_501_NOT_IMPLEMENTED:
		status = stat.StatusCode_ERR_NOT_IMPLEMENTED
	case COAP_503_SERVICE_UNAVAILABLE:
		status = stat.StatusCode_ERR_SERVICE_UNAVAILABLE
	default:
		status = stat.StatusCode_ERR_INTERNAL_SERVER_ERROR
	}

	return status
}

//Get SENML Json Value
func GetSENMLJsonValue(typeMap map[ID]string, rec senml.Record) (*string, error) {
	var val string

	//If record name is not nil
	if rec.Name != nil {
		//convert resourceID to integer
		rid, err := strconv.Atoi(*rec.Name)
		if err != nil {
			return nil, err
		}

		//Read the value
		dType := typeMap[ID(rid)]
		switch dType {
		case "String":
			if rec.StringValue != nil {
				val = fmt.Sprint(*rec.StringValue)
			}

		case "Integer":
			if rec.DataValue != nil {
				val = fmt.Sprint(*rec.DataValue)
			}
		case "Float":
			if rec.Value != nil {
				val = fmt.Sprintf("%f", *rec.Value)
			}
		case "bool":
			if rec.BoolValue != nil {
				val = strconv.FormatBool(*rec.BoolValue)
			}
		default:
			//TODO:: Check what do here
			return nil, errors.New("unknown Data Type for ResourceID")
		}
	}

	return &val, nil

}

// Preparing config data from the viper.
func GetUpdatedCfgData(token uint32) (*string, error) {

	var cfgdata string

	// Read all the config data as key, interface{} map.
	c := viper.AllSettings()

	// Marshal configs into supported config types */
	switch cfg.Config.DevConfigType {
	case "json":
		b, err := json.MarshalIndent(c, "", "  ")
		if err != nil {
			return nil, err
		}
		cfgdata = string(b)

	case "toml":
		t, err := toml.TreeFromMap(c)
		if err != nil {
			return nil, err
		}
		cfgdata = t.String()

	default:
		return nil, errors.New("unknown config type")
	}
	return &cfgdata, nil

}

// Create a Device Config from the SENML-JSON
func createDevCfg(token uint32, cfgJson *senml.Msg, cfgData string) (*string, error) {

	// Setting viper to read and modify the config data.
	viper.SetConfigType(cfg.Config.DevConfigType)             // config type
	err := viper.ReadConfig(bytes.NewBuffer([]byte(cfgData))) // Find and read the config file
	if err != nil {                                           // Handle errors reading the config file
		log.Errorf("LwM2M::Error parsing config data: %s. Error:: %s \n", cfgData, err)
		return nil, err
	}
	log.Debugf("LwM2M::Reading config data for ObjectID %s, InstanceID %s", viper.GetString("Info.ObjectID"), viper.GetString("Info.InstanceID"))

	//Get resource ID to data type map.
	rIDToTypeMap, err := GetResourceIDDataType(viper.GetString("Info.ObjectID"))
	if err != nil {
		log.Errorf("LwM2M::Error getting ResourceID to data type map:: Error %s", err)
		return nil, err
	}

	//Update Config data now
	records := cfgJson.Records
	for _, rec := range records {
		val, err := GetSENMLJsonValue(rIDToTypeMap, rec)
		if (err != nil) || (len(*val) == 0) {
			continue
		}
		cfgName := *rec.Name
		viper.Set(string("Config.")+cfgName, *val)
		log.Debugf("LwM2M::Updated resource id %+v with Value %+v", cfgName, *val)
		log.Debugf("LwM2M::Reading back ObjectID %s, InstanceID %s: Resource %s: Value  %s", viper.GetString("Info.ObjectID"), viper.GetString("Info.InstanceID"), cfgName, viper.GetString(string("Config.")+cfgName))
	}

	rCfg, err := GetUpdatedCfgData(token)
	if err != nil {
		log.Errorf("LwM2M::Error unable to prepare config data with the response and request config. Error %s", err.Error())
		return nil, err
	}

	viper.Reset()

	return rCfg, nil
}

// Function Decode content from  SWNMLJSON to toml
func decodeContent(ct Lwm2mContentType, data string) (*senml.Msg, error) {
	var content senml.Msg
	var err error

	bdata := []byte(data)

	switch ct {
	case LWM2M_CONTENT_JSON:
		content, err = senml.DecodeMsg(bdata, senml.JSON)
		if err != nil {
			log.Debugf("LwM2M::Failed to decode senml message. Error:: %s", err.Error())
		}

		//DEBUG Purpose
		for i, rec := range content.Records {
			log.Tracef("LwM2M::Records[%d] is %+v", i, rec)
		}

	default:
		log.Debugf("LwM2M::Unsupported content type in response message.")
	}

	return &content, err
}

// Process LwM2M server read data. this requires copying of content
func processReadRespData(resp *ResponseMsg, cfgData string) (*string, error) {

	//Decode ContentType
	cfgJson, err := decodeContent(Lwm2mContentType(resp.Format), resp.Message)
	if err != nil {
		log.Errorf("LwM2M::Problem while decoding SENML-JSON Reply from LWM2M server.")
		return nil, err
	}

	// Create new cfgData for the read config
	rCfgData, err := createDevCfg(resp.Token, cfgJson, cfgData)
	if err != nil {
		log.Errorf("LwM2M::Failed to create config file for the read config. Error %+v", err)
		return nil, err
	}

	return rCfgData, nil

}

// Read Config request
func ReadConfig(r *spec.Lwm2MConfigReqMsg, newCfg *string) stat.StatusCode {
	var ret stat.StatusCode

	// Get URL
	urlext := URItoString(r.Uri, Read)
	log.Infof("LwM2M::Request from controller:: Read config message for %s URI %s", r.Device.Uuid, *urlext)

	//Prepare message
	msg := PrepareMsgForLwm2mServer(Read, &r.Device.Uuid, urlext, nil)

	//Send and wait for response.
	resp, err := transmit(r.Device.Uuid, *msg)
	if err == nil {

		// handle response
		if resp.Status == COAP_205_CONTENT {

			// Decode the responded config data from LwM2M server.
			cfg, err := processReadRespData(resp, r.CfgData)
			if (err == nil) && (newCfg != nil) {
				*newCfg = *cfg
				ret = stat.StatusCode_STATUS_CONTENT

			} else {
				log.Errorf("LwM2M::Failed to process read response from  LwM2M server. Error :: %+v", err)
				ret = stat.StatusCode_ERR_INTERNAL_SERVER_ERROR
			}

		} else {
			/* if not content */
			ret = LwM2MResponseCode(resp.Status)
		}

	} else {
		log.Errorf("LwM2M::Failed to connect to LwM2M server. Error :: %+v", err)
		ret = stat.StatusCode_ERR_FAILED_TO_CONNECT_SERVICE
	}
	log.Infof("LwM2M::Responding to Controller:: Read config message for %s URI %s with status %s", r.Device.Uuid, *urlext, stat.StatusCode_name[int32(ret)])
	return ret
}

// Start processing fro write request
func GetValueFromCfgData(cfgData string, rId ID) (*string, error) {

	// Setting viper to read  the config data.
	viper.SetConfigType(cfg.Config.DevConfigType)             // config type
	err := viper.ReadConfig(bytes.NewBuffer([]byte(cfgData))) // Find and read the config file
	if err != nil {                                           // Handle errors reading the config file
		log.Errorf("LwM2M::Error parsing config data: %s. Error:: %s \n", cfgData, err)
		return nil, err
	}
	log.Debugf("LwM2M::Reading ObjectID %s, InstanceID %s", viper.GetString("Info.ObjectID"), viper.GetString("Info.InstanceID"))

	val := viper.GetString(string("Config.") + strconv.Itoa(int(rId)))
	if val == "" {
		log.Errorf("LwM2M::No Value read for %d from config data provided.", int(rId))
		return nil, errors.New("no config data available for resource")
	}

	log.Debugf("LwM2M::Configuration data read is %s for Resource ID %d", val, rId)

	return &val, nil
}

//Write Config request
func WriteConfig(r *spec.Lwm2MConfigReqMsg) stat.StatusCode {

	ret := stat.StatusCode_ERR_NOT_FOUND

	log.Infof("LwM2M::Request from controller:: Write config message for %s URI /%d/%d/%d", r.Device.Uuid, r.Uri.Object, r.Uri.Instance, r.Uri.Resource)

	// Get resource list to be updated
	rIds, err := GetResourcesIdByOperation(strconv.Itoa(int(r.Uri.Object)), WRITE)
	if err != nil {
		log.Errorf("LwM2M::Failed to find the resource id for %s in schema. Error :: %+v", WRITE, err)
		ret = stat.StatusCode_ERR_NOT_FOUND
	}

	// Iterating through all writable resources.
	for _, rId := range rIds {
		uri := "/" + strconv.Itoa(int(r.Uri.Object)) + "/" + strconv.Itoa(int(r.Uri.Instance)) + "/" + strconv.Itoa(int(rId))

		// Get the value for the resourceID to write from toml.
		val, err := GetValueFromCfgData(r.CfgData, rId)
		if err != nil {
			log.Errorf("LwM2M::Failed to get resource id %d value from the config data. Error: %s", rId, err)
			return stat.StatusCode_ERR_INTERNAL_SERVER_ERROR
		}

		if val == nil {
			log.Errorf("LwM2M::Failed to find resource id %d or it's value is nil", rId)
			return stat.StatusCode_ERR_INTERNAL_SERVER_ERROR
		}

		//Prepare message
		msg := PrepareMsgForLwm2mServer(Write, &r.Device.Uuid, &uri, val)

		//Send and wait for response.
		resp, err := transmit(r.Device.Uuid, *msg)
		if err == nil {
			// handle response
			if resp.Status != COAP_204_CHANGED {
				log.Errorf("LwM2M::Failed to update config for %s object %d instance %d and resource %d value %s. Error:: %s",
					r.Device.Uuid, r.Uri.Object, r.Uri.Instance, rId, *val, stat.StatusCode_name[int32(resp.Status)])
				break
			}
			ret = LwM2MResponseCode(resp.Status)
		} else {
			log.Errorf("LwM2M::Failed to connect to LwM2M server. Error :: %+v", err)
			ret = stat.StatusCode_ERR_FAILED_TO_CONNECT_SERVICE
		}

	}
	log.Infof("LwM2M::Responding to Controller:: Write config message for %s URI /%d/%d/%d with status %s", r.Device.Uuid, r.Uri.Object,
		r.Uri.Instance, r.Uri.Resource, stat.StatusCode_name[int32(ret)])
	return ret
}

// Execute  Command request
func ExecCommand(r *spec.Lwm2MConfigReqMsg) stat.StatusCode {
	log.Debugf("LwM2M::Execute Config request %+v", r)
	ret := stat.StatusCode_ERR_NOT_FOUND

	log.Infof("LwM2M::Request from Controller:: Received execute message for %s URI /%d/%d/%d", r.Device.Uuid, r.Uri.Object, r.Uri.Instance, r.Uri.Resource)

	// Get resource list to be updated
	rIds, err := GetResourcesIdByOperation(strconv.Itoa(int(r.Uri.Object)), EXECUTE)
	if err != nil {
		log.Errorf("LwM2M::Failed to find the resource id for %s in schema. Error :: %+v", EXECUTE, err)
		ret = stat.StatusCode_ERR_NOT_FOUND
	}

	// Iterating through all writable resources.
	for _, rId := range rIds {
		ret = stat.StatusCode_ERR_NOT_FOUND

		if int32(rId) == r.Uri.Resource {
			uri := "/" + strconv.Itoa(int(r.Uri.Object)) + "/" + strconv.Itoa(int(r.Uri.Instance)) + "/" + strconv.Itoa(int(rId))

			//Prepare message
			msg := PrepareMsgForLwm2mServer(Execute, &r.Device.Uuid, &uri, nil)

			//Send and wait for response.
			resp, err := transmit(r.Device.Uuid, *msg)
			if err == nil {
				// handle response
				if resp.Status != COAP_204_CHANGED {
					log.Errorf("LwM2M::Failed to execute %s object %d instance %d and resource %d. Error:: %s",
						r.Device.Uuid, r.Uri.Object, r.Uri.Instance, rId, stat.StatusCode_name[int32(resp.Status)])
					break
				}
				ret = LwM2MResponseCode(resp.Status)
			} else {
				log.Errorf("LwM2M::Failed to connect to LwM2M server. Error :: %+v", err)
				ret = stat.StatusCode_ERR_FAILED_TO_CONNECT_SERVICE
			}
			break
		}

	}
	log.Infof("LwM2M::Responding to Controller:: execute message for %s URI /%d/%d/%d with status %s", r.Device.Uuid, r.Uri.Object, r.Uri.Instance,
		r.Uri.Resource, stat.StatusCode_name[int32(ret)])
	return ret
}

// Enable notification for objects
func ObserveConfig(r *spec.Lwm2MConfigReqMsg) stat.StatusCode {
	log.Debugf("LwM2M::Enabling Observation on config request: %+v", r)
	var ret stat.StatusCode

	// Get URL
	urlext := URItoString(r.Uri, Observe)

	log.Infof("LwM2M::Request from Controller:: Received Observe message for %s URI %s", r.Device.Uuid, *urlext)

	//Prepare message
	msg := PrepareMsgForLwm2mServer(Observe, &r.Device.Uuid, urlext, nil)

	//Send and wait for response.
	resp, err := transmit(r.Device.Uuid, *msg)
	if err == nil {
		// Handle response
		if resp.Status != COAP_204_CHANGED {
			log.Debugf("LwM2M::Failed to enable notification on  %s device %s resource. Error:: %s",
				r.Device.Uuid, *urlext, stat.StatusCode_name[int32(resp.Status)])
		}
		ret = LwM2MResponseCode(resp.Status)
	} else {
		log.Errorf("LwM2M::Failed to connect to LwM2M server. Error :: %+v", err)
		ret = stat.StatusCode_ERR_FAILED_TO_CONNECT_SERVICE
	}
	log.Infof("LwM2M::Responding to Controller:: Observe message for %s URI %s with status %s", r.Device.Uuid, *urlext, stat.StatusCode_name[int32(ret)])
	return ret
}

// Disable notification for objects
// As of this function is similar to enable it it remains same
// we can abstract this function to handle both
func CancelObservationOnConfig(r *spec.Lwm2MConfigReqMsg) stat.StatusCode {
	log.Debugf("LwM2M::Cancel observation on config request %+v", r)
	var ret stat.StatusCode

	// Get URL
	urlext := URItoString(r.Uri, Cancel)
	log.Infof("LwM2M::Request from Controller:: Received Cancel Observe message for %s URI %s", r.Device.Uuid, *urlext)
	//Prepare message
	msg := PrepareMsgForLwm2mServer(Cancel, &r.Device.Uuid, urlext, nil)

	//Send and wait for response.
	resp, err := transmit(r.Device.Uuid, *msg)
	if err == nil {
		// Handle response
		if resp.Status != COAP_NO_ERROR {
			log.Debugf("LwM2M::Failed to disable notification on  %s device %s resource. Error:: %s",
				r.Device.Uuid, *urlext, stat.StatusCode_name[int32(resp.Status)])
		}
		ret = LwM2MResponseCode(resp.Status)
	} else {
		log.Errorf("LwM2M::Failed to connect to LwM2M server. Error :: %+s", err.Error())
		ret = stat.StatusCode_ERR_FAILED_TO_CONNECT_SERVICE
	}
	log.Infof("LwM2M::Responding to Controller:: Cancel Observe message for %s URI %s with status %s", r.Device.Uuid, *urlext, stat.StatusCode_name[int32(ret)])
	return ret
}
