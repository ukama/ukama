package senml_test

import (
	"fmt"
	"testing"

	"lwm2m-gateway/pkg/senml"
)

//const lwdata string = `{"bn":"/3303/0/","e":[{"n":"5700","v":45000},{"n":"5601","v":0},{"n":"5602","v":45000},{"n":"5603","v":0},{"n":"5604","v":0},{"n":"5701","sv":"milliCelsius"},{"n":"5750","sv":"Pmic"}]}`
const lwdata string = `{"bn":"/3303/0/","e":[{"n":"5700","v":45000},{"n":"5601","v":0},{"n":"5602","v":45000},{"n":"5603","v":0},{"n":"5604","v":0},{"n":"5701","vs":"milliCelsius"},{"n":"5750","vs":"Pmic"}]}`
const lwnew string = `  {"bn":"/3328/0/5821/","e":[{"n":"5821","v":1700}]}`

//const lwdata string = `{"bn":"/3303/0/"}`

// https://tools.ietf.org/html/rfc8428#section-7
const xmlData string = `<sensml xmlns="urn:ietf:params:xml:ns:senml">
	<senml bn="urn:dev:ow:10e2073a0108006:" bt="1.276020076001e+09"
	bu="A" bver="5" n="voltage" u="V" v="120.1"></senml>
	<senml n="current" t="-5" v="1.2"></senml>
	<senml n="current" t="-4" v="1.3"></senml>
	<senml n="current" t="-3" v="1.4"></senml>
	<senml n="current" t="-2" v="1.5"></senml>
	<senml n="current" t="-1" v="1.6"></senml>
	<senml n="current" v="1.7"></senml>
  </sensml>`

func TestNewDecodeJSON(t *testing.T) {
	data, err := senml.DecodeMsg([]byte(lwnew), senml.JSON)
	if err != nil {
		t.Error("Decoding JSON failed: ", err)
		return
	}
	fmt.Printf(" Json data is data %+v, senml.JSON %+v\n", data, senml.JSON)

	records := data.Records
	for i, rec := range records {
		fmt.Printf("Record [%d] %+v.\n", i, rec)
		fmt.Printf("Record is:: %+v,   %+v ,  %+v ,  %+v,  %+v ", rec.BaseName, rec.BaseSum, rec.BaseTime, rec.BaseUnit, rec.BaseValue)
		fmt.Printf("Record is:: %+v,   %+v ,  %+v ,  %+v,  %+v, %+v", rec.BaseVersion, rec.BoolValue, rec.DataValue, rec.Name, rec.StringValue, rec.Sum)
		fmt.Printf("Record is:: %+v,   %+v ,  %+v ,  %+v,  %+v, %+v", rec.Time, rec.UpdateTime, rec.Value, rec.XMLName, rec.XMLName.Local, rec.XMLName.Space)
	}
}

func TestDecodeJSON(t *testing.T) {
	data, err := senml.DecodeMsg([]byte(lwdata), senml.JSON)
	if err != nil {
		t.Error("Decoding JSON failed: ", err)
		return
	}
	fmt.Printf(" Json data is data %+v, senml.JSON %+v\n", data, senml.JSON)

	records := data.Records
	for i, rec := range records {
		fmt.Printf("Record [%d] %+v.\n", i, rec)
		fmt.Printf("Record is:: %+v,   %+v ,  %+v ,  %+v,  %+v ", rec.BaseName, rec.BaseSum, rec.BaseTime, rec.BaseUnit, rec.BaseValue)
		fmt.Printf("Record is:: %+v,   %+v ,  %+v ,  %+v,  %+v, %+v", rec.BaseVersion, rec.BoolValue, rec.DataValue, rec.Name, rec.StringValue, rec.Sum)
		fmt.Printf("Record is:: %+v,   %+v ,  %+v ,  %+v,  %+v, %+v", rec.Time, rec.UpdateTime, rec.Value, rec.XMLName, rec.XMLName.Local, rec.XMLName.Space)
	}
}

func TestDecodeXML(t *testing.T) {
	_, err := senml.Decode([]byte(xmlData), senml.XML)
	if err != nil {
		t.Error("Decoding XML failed: ", err)
		return
	}
}

func TestDecodeInvalidFormat(t *testing.T) {
	_, err := senml.Decode(nil, -1)
	if err == nil {
		t.Error("Decoding an invalid format should result in an error")
	}
}

func TestEncodeJSON(t *testing.T) {
	message, err := senml.DecodeMsg([]byte(lwdata), senml.JSON)
	if err != nil {
		t.Error("Decoding JSON failed: ", err)
		return
	}
	fmt.Printf(" Json data is data %+v, senml.JSON %+v\n", message, senml.JSON)
	encodedMessage, err := message.Encode(senml.JSON)
	if err != nil {
		t.Error("Encoding message to JSON failed: ", err)
		return
	}
	fmt.Printf(" Json data is message %+v, encodedMessage %+v, senml.JSON %+v", message, string(encodedMessage), senml.JSON)

	if len(encodedMessage) == 0 {
		t.Error("Encoding to JSON resulted in an empty message")
	}
}

func TestEncodeXML(t *testing.T) {
	message, err := senml.Decode([]byte(xmlData), senml.XML)
	if err != nil {
		t.Error("Decoding XML failed: ", err)
		return
	}

	encodedMessage, err := message.Encode(senml.XML)
	if err != nil {
		t.Error("Encoding message to XML failed: ", err)
		return
	}
	if len(encodedMessage) == 0 {
		t.Error("Encoding to XML resulted in an empty message")
	}
}

func TestEncodeInvalidFormat(t *testing.T) {
	message := senml.Message{}
	_, err := message.Encode(-1)
	if err == nil {
		t.Error("Encoding message to an invalid format should result in an error")
		return
	}
}

func TestResolveUnsupportedSenMLVersion(t *testing.T) {
	var unsupportedVersion = 11
	var name = "test"
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseVersion: &unsupportedVersion,
				Name:        &name,
				Value:       &value,
			},
		},
	}

	_, err := message.Resolve()
	if err == nil {
		t.Error("Resolving an unsupported SenML version should result in an error")
		return
	}
}

func TestResolveBaseVersionIsSetIfLowerThanMaximumSupported(t *testing.T) {
	var lowerVersion = 5
	var name = "test"
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseVersion: &lowerVersion,
				BaseName:    &name,
				Value:       &value,
			},
			{
				Value: &value,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving a lower SenML version than the maximum supported version failed:", err)
		return
	}

	for _, record := range resolvedMessage.Records {
		if record.BaseVersion == nil {
			t.Error("The BaseVersion attribute is not set if the version is lower than the maximum supported version")
			return
		}
		if *record.BaseVersion != lowerVersion {
			t.Error("The BaseVersion attribute is not set to the BaseVersion in the unresolved message")
			return
		}
	}
}

func TestResolveRecordsHaveDifferentVersion(t *testing.T) {
	var version = 5
	var differentVersion = 6
	var name = "test"
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseVersion: &version,
				BaseName:    &name,
				Value:       &value,
			},
			{
				BaseVersion: &differentVersion,
				Value:       &value,
			},
		},
	}

	_, err := message.Resolve()
	if err == nil {
		t.Error("Resolving a SenML message which contains records with different version should result in an error")
		return
	}
}

func TestResolveNameContainsInvalidSymbols(t *testing.T) {
	var name = "test("
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:  &name,
				Value: &value,
			},
		},
	}

	_, err := message.Resolve()
	if err == nil {
		t.Error("Resolving a record with a name which contains invalid symbols should result in an error")
		return
	}
}

func TestResolveNameStartsWithInvalidSymbols(t *testing.T) {
	var name = "-test"
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:  &name,
				Value: &value,
			},
		},
	}

	_, err := message.Resolve()
	if err == nil {
		t.Error("Resolving a record with a name which starts with an invalid symbol should result in an error")
		return
	}
}

func TestResolveNoName(t *testing.T) {
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				Value: &value,
			},
		},
	}

	_, err := message.Resolve()
	if err == nil {
		t.Error("Resolving a record with no name should result in an error")
		return
	}
}

func TestResolveNoValue(t *testing.T) {
	var name = "test"
	message := senml.Message{
		Records: []senml.Record{
			{
				Name: &name,
			},
		},
	}

	_, err := message.Resolve()
	if err == nil {
		t.Error("Resolving a record with no value or sum should result in an error")
		return
	}
}

func TestResolveValue(t *testing.T) {
	var name = "test"
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:  &name,
				Value: &value,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving a record with a value should not result in an error", err)
		return
	}

	if resolvedMessage.Records[0].Value == nil {
		t.Error("The record in the resolved message has no value")
		return
	}

	if *resolvedMessage.Records[0].Value != value {
		t.Error("The value field has a different value than expected")
		return
	}
}

func TestResolveBoolValue(t *testing.T) {
	var name = "test"
	var boolValue bool = true
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:      &name,
				BoolValue: &boolValue,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving a record with a bool value should not result in an error", err)
		return
	}

	if resolvedMessage.Records[0].BoolValue == nil {
		t.Error("The record in the resolved message has no bool value")
		return
	}

	if *resolvedMessage.Records[0].BoolValue != boolValue {
		t.Error("The bool value field has a different value than expected")
		return
	}
}

func TestResolveStringValue(t *testing.T) {
	var name = "test"
	var stringValue string = "value"
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:        &name,
				StringValue: &stringValue,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving a record with a string value should not result in an error", err)
		return
	}

	if resolvedMessage.Records[0].StringValue == nil {
		t.Error("The record in the resolved message has no string value")
		return
	}

	if *resolvedMessage.Records[0].StringValue != stringValue {
		t.Error("The string value field has a different value than expected")
		return
	}
}

func TestResolveDataValue(t *testing.T) {
	var name = "test"
	var dataValue string = "data"
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:      &name,
				DataValue: &dataValue,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving a record with a data value should not result in an error", err)
		return
	}

	if resolvedMessage.Records[0].DataValue == nil {
		t.Error("The record in the resolved message has no data value")
		return
	}

	if *resolvedMessage.Records[0].DataValue != dataValue {
		t.Error("The data value field has a different value than expected")
		return
	}
}

func TestResolveSum(t *testing.T) {
	var name = "test"
	var sum float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				Name: &name,
				Sum:  &sum,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving a record with a sum should not result in an error", err)
		return
	}

	if resolvedMessage.Records[0].Sum == nil {
		t.Error("The record in the resolved message has no sum")
		return
	}

	if *resolvedMessage.Records[0].Sum != sum {
		t.Error("The sum field has a different value than expected")
		return
	}
}

func TestResolveUnit(t *testing.T) {
	var name = "test"
	var value float64 = 1
	var unit = "unit"
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:  &name,
				Value: &value,
				Unit:  &unit,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].Unit == nil {
		t.Error("The record in the resolved message has no unit")
		return
	}

	if *resolvedMessage.Records[0].Unit != unit {
		t.Error("The unit field has a different value than expected")
		return
	}
}

func TestResolveUpdateTime(t *testing.T) {
	var name = "test"
	var value float64 = 1
	var updateTime float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:       &name,
				Value:      &value,
				UpdateTime: &updateTime,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].UpdateTime == nil {
		t.Error("The record in the resolved message has no update time")
		return
	}

	if *resolvedMessage.Records[0].UpdateTime != updateTime {
		t.Error("The update time field has a different value than expected")
		return
	}
}

func TestResolveRelativeToAbsoluteTime(t *testing.T) {
	var name = "test"
	var value float64 = 1
	var time float64 = 2
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:  &name,
				Value: &value,
				Time:  &time,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].Time == nil {
		t.Error("The record in the resolved message has no time")
		return
	}

	if *resolvedMessage.Records[0].Time == time {
		t.Error("The time field was not resolved")
	}
}

func TestResolveAbsoluteTime(t *testing.T) {
	var name = "test"
	var value float64 = 1
	var time float64 = 2 ^ 28
	message := senml.Message{
		Records: []senml.Record{
			{
				Name:  &name,
				Value: &value,
				Time:  &time,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].Time == nil {
		t.Error("The record in the resolved message has no time")
		return
	}

	if *resolvedMessage.Records[0].Time != time {
		t.Error("The value of the time field was changed, but the RFC specifies that values of 2^28 or over should not be changed")
	}
}

func TestResolveOrderIsChronological(t *testing.T) {
	var baseName = "test"
	var value float64 = 1
	var value2 float64 = 2
	var value3 float64 = 3
	var value4 float64 = 4
	var time3 float64 = 3
	var time4 float64 = 4
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseName: &baseName,
				Value:    &value4,
				Time:     &time4,
			},
			{
				Value: &value,
			},
			{
				Value: &value2,
			},
			{
				Value: &value3,
				Time:  &time3,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].Value == nil {
		t.Error("The record in the resolved message has no value")
		return
	}

	if *resolvedMessage.Records[0].Value != value {
		t.Error("The records are not in chronological order")
		return
	}

	if resolvedMessage.Records[1].Value == nil {
		t.Error("The record in the resolved message has no value")
		return
	}

	if *resolvedMessage.Records[1].Value != value2 {
		t.Error("The records are not in chronological order")
		return
	}

	if resolvedMessage.Records[2].Value == nil {
		t.Error("The record in the resolved message has no value")
		return
	}

	if *resolvedMessage.Records[2].Value != value3 {
		t.Error("The records are not in chronological order")
		return
	}

	if resolvedMessage.Records[3].Value == nil {
		t.Error("The record in the resolved message has no value")
		return
	}

	if *resolvedMessage.Records[3].Value != value4 {
		t.Error("The records are not in chronological order")
		return
	}
}

func TestResolveBaseName(t *testing.T) {
	var baseName = "base/"
	var name = "test"
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseName: &baseName,
				Value:    &value,
			},
			{
				Name:  &name,
				Value: &value,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].BaseName != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[0].Name == nil {
		t.Error("The record in the resolved message has no name")
		return
	}

	if *resolvedMessage.Records[0].Name != baseName {
		t.Error("The base attribute was not properly concatenated with the field")
		return
	}

	if resolvedMessage.Records[1].BaseName != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[1].Name == nil {
		t.Error("The record in the resolved message has no name")
		return
	}

	if *resolvedMessage.Records[1].Name != baseName+name {
		t.Error("The base attribute was not properly concatenated with the field")
		return
	}
}

func TestResolveBaseTime(t *testing.T) {
	var baseTime float64 = 2 ^ 28
	var baseName = "test"
	var baseValue float64 = 1
	var time float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseTime:  &baseTime,
				BaseName:  &baseName,
				BaseValue: &baseValue,
			},
			{
				Time: &time,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].BaseTime != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[0].Time == nil {
		t.Error("The record in the resolved message has no time")
		return
	}

	if *resolvedMessage.Records[0].Time != baseTime {
		t.Error("The base attribute was not properly added to the field")
		return
	}

	if resolvedMessage.Records[1].BaseTime != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[1].Time == nil {
		t.Error("The record in the resolved message has no time")
		return
	}

	if *resolvedMessage.Records[1].Time != baseTime+time {
		t.Error("The base attribute was not properly added to the field")
		return
	}
}

func TestResolveBaseUnit(t *testing.T) {
	var baseUnit = "bu"
	var baseName = "test"
	var baseValue float64 = 1
	var unit string = "u"
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseUnit:  &baseUnit,
				BaseName:  &baseName,
				BaseValue: &baseValue,
			},
			{
				Unit: &unit,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].BaseUnit != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[0].Unit == nil {
		t.Error("The record in the resolved message has no unit")
		return
	}

	if *resolvedMessage.Records[0].Unit != baseUnit {
		t.Error("The base attribute was not properly set")
		return
	}

	if resolvedMessage.Records[1].BaseUnit != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[1].Unit == nil {
		t.Error("The record in the resolved message has no unit")
		return
	}

	if *resolvedMessage.Records[1].Unit != unit {
		t.Error("The field was replaced with the base attribute")
		return
	}
}

func TestResolveBaseValue(t *testing.T) {
	var baseValue float64 = 1
	var baseName = "test"
	var value float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseValue: &baseValue,
				BaseName:  &baseName,
			},
			{
				Value: &value,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].BaseValue != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[0].Value == nil {
		t.Error("The record in the resolved message has no value")
		return
	}

	if *resolvedMessage.Records[0].Value != baseValue {
		t.Error("The base attribute was not properly added to the field")
		return
	}

	if resolvedMessage.Records[1].BaseValue != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[1].Value == nil {
		t.Error("The record in the resolved message has no value")
		return
	}

	if *resolvedMessage.Records[1].Value != baseValue+value {
		t.Error("The base attribute was not properly added to the field")
		return
	}
}

func TestResolveBaseSum(t *testing.T) {
	var baseSum float64 = 1
	var baseName = "test"
	var sum float64 = 1
	message := senml.Message{
		Records: []senml.Record{
			{
				BaseSum:  &baseSum,
				BaseName: &baseName,
			},
			{
				Sum: &sum,
			},
		},
	}

	resolvedMessage, err := message.Resolve()
	if err != nil {
		t.Error("Resolving the record failed", err)
		return
	}

	if resolvedMessage.Records[0].BaseSum != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[0].Sum == nil {
		t.Error("The record in the resolved message has no sum")
		return
	}

	if *resolvedMessage.Records[0].Sum != baseSum {
		t.Error("The base attribute was not properly added to the field")
		return
	}

	if resolvedMessage.Records[1].BaseSum != nil {
		t.Error("The resolved record has a base attribute set")
		return
	}

	if resolvedMessage.Records[1].Sum == nil {
		t.Error("The record in the resolved message has no sum")
		return
	}

	if *resolvedMessage.Records[1].Sum != baseSum+sum {
		t.Error("The base attribute was not properly added to the field")
		return
	}
}

func TestInvalidNameErrorFirstCharacterInvalid(t *testing.T) {
	err := &senml.InvalidNameError{
		Reason: senml.FirstCharacterInvalid,
	}
	message := err.Error()
	if message == "" {
		t.Error("The error message is empty.")
	}
}

func TestInvalidNameErrorContainsInvalidCharacter(t *testing.T) {
	err := &senml.InvalidNameError{
		Reason: senml.ContainsInvalidCharacter,
	}
	message := err.Error()
	if message == "" {
		t.Error("The error message is empty.")
	}
}

func TestInvalidNameErrorEmpty(t *testing.T) {
	err := &senml.InvalidNameError{
		Reason: senml.Empty,
	}
	message := err.Error()
	if message == "" {
		t.Error("The error message is empty.")
	}
}

func TestInvalidNameErrorUnknown(t *testing.T) {
	err := &senml.InvalidNameError{
		Reason: -1,
	}
	message := err.Error()
	if message == "" {
		t.Error("The error message is empty.")
	}
}

func TestUnsupportedVersionError(t *testing.T) {
	err := &senml.UnsupportedVersionError{
		SupportedVersion: 10,
		GivenVersion:     11,
	}
	message := err.Error()
	if message == "" {
		t.Error("The error message is empty.")
	}
}

func TestDifferentVersionError(t *testing.T) {
	err := &senml.DifferentVersionError{
		CurrentVersion: 10,
		GivenVersion:   11,
	}
	message := err.Error()
	if message == "" {
		t.Error("The error message is empty.")
	}
}

func TestMissingValueError(t *testing.T) {
	err := &senml.MissingValueError{}
	message := err.Error()
	if message == "" {
		t.Error("The error message is empty.")
	}
}

func TestUnsupportedFormatError(t *testing.T) {
	err := &senml.UnsupportedFormatError{
		GivenFormat: -1,
	}
	message := err.Error()
	if message == "" {
		t.Error("The error message is empty.")
	}
}
