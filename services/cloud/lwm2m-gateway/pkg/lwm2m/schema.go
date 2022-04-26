package lwm2m

import (
	"errors"
	"io/ioutil"

	cfg "lwm2m-gateway/pkg/config"
	"strings"

	toml "github.com/pelletier/go-toml"

	log "github.com/sirupsen/logrus"
)

type item struct {
	Name              string
	ResourceID        ID
	Operations        string
	MultipleInstances string
	Mandatory         string
	Type              string
	RangeEnumeration  string
	Units             string
	Description       string
}

type resource struct {
	Item []item `toml:"Item"`
}

type object struct {
	Name              string
	Description1      string
	ObjectID          ID
	ObjectURN         string
	LWM2MVersion      int16
	ObjectVersion     string
	MultipleInstances string
	Mandatory         string
	Description2      string
	Resource          resource `toml:"Resources"`
}

type lwm2m struct {
	Object object `toml:"Object"`
}

type schema struct {
	Lwm2m lwm2m `toml:"LWM2M"`
}

// Read toml file for the Object Schema
func readSchemaFromFile(fileName string) ([]byte, error) {

	// Get complete path
	fpath := cfg.Config.SchemaDir + fileName + SCHEMA_POSTFIX

	//Read file
	fdata, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Debugf("Schema::Error reading file %s Error %+v", fpath, err)
		return []byte(""), err
	}
	return fdata, nil
}

// Unmarshal the Object schema toml file.
func unmarshalSchema(fileName string) (*schema, error) {

	schema := schema{
		Lwm2m: lwm2m{},
	}

	//Read schema
	fdata, err := readSchemaFromFile(fileName)
	if err != nil {
		log.Errorf("Schema::Failed reading file %s Error:: %+v", fileName, err)
		return nil, err
	}

	//Unmarshal
	err = toml.Unmarshal(fdata, &schema)
	if err != nil {
		log.Errorf("Schema::Failed to parse file %s Error:: %+v", fileName, err)
		return nil, err
	}

	return &schema, nil
}

// Get all the resource id of the Object by type od operation it support.
func GetResourcesIdByOperation(fileName string, op string) ([]ID, error) {

	var WrtResourceID []ID
	//Read schema
	schema, err := unmarshalSchema(fileName)
	if err != nil {
		return WrtResourceID, err
	}

	// Add all writable resources to the list
	log.Debugf("Schema::Resources present in Object %s are %d. Looking for resources who support %s operation", fileName, len(schema.Lwm2m.Object.Resource.Item), op)
	for _, item := range schema.Lwm2m.Object.Resource.Item {

		if strings.Contains(item.Operations, op) {
			WrtResourceID = append(WrtResourceID, item.ResourceID)
		}
	}
	log.Debugf("Schema::ResourceID supporting %s Operation in object %s are %+v", op, fileName, WrtResourceID)
	return WrtResourceID, nil
}

// Check if the operation is supported by resource Id.
func CheckIfOperationAvailOnResourcesId(fileName string, resourceId ID, op string) (bool, error) {

	opSupported := false

	// Read ResourceID by operatoion type supported.
	WrtResourceID, err := GetResourcesIdByOperation(fileName, op)
	if err != nil {
		return opSupported, err
	}

	// Iterate over the range returned.
	for _, id := range WrtResourceID {
		if id == resourceId {
			opSupported = true
			break
		}
	}

	return opSupported, nil
}

// Read the data type for value stored against resourceID in config file.
func GetResourceIDDataType(fileName string) (map[ID]string, error) {

	if fileName == "" {
		return nil, errors.New("invalid ObjectId")
	}

	typeMap := make(map[ID]string)

	//Read schema
	schema, err := unmarshalSchema(fileName)
	if err != nil {
		return nil, err
	}

	// Add all writable resources to the list
	log.Debugf("Schema::Resources present in Object %s are %d.", fileName, len(schema.Lwm2m.Object.Resource.Item))
	for _, item := range schema.Lwm2m.Object.Resource.Item {
		typeMap[item.ResourceID] = item.Type
	}
	log.Debugf("Schema::ResourceID and Type map for object %s is %+v", fileName, typeMap)

	return typeMap, nil
}
