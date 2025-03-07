// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: subscriber.proto

package gen

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/mwitkow/go-proto-validators"
	_ "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	regexp "regexp"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *ListSubscribersRequest) Validate() error {
	return nil
}
func (this *ListSubscribersResponse) Validate() error {
	for _, item := range this.Subscribers {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Subscribers", err)
			}
		}
	}
	return nil
}

var _regex_DeleteSubscriberRequest_SubscriberId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *DeleteSubscriberRequest) Validate() error {
	if !_regex_DeleteSubscriberRequest_SubscriberId.MatchString(this.SubscriberId) {
		return github_com_mwitkow_go_proto_validators.FieldError("SubscriberId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.SubscriberId))
	}
	if this.SubscriberId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("SubscriberId", fmt.Errorf(`value '%v' must not be an empty string`, this.SubscriberId))
	}
	return nil
}

var _regex_GetByNetworkRequest_NetworkId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *GetByNetworkRequest) Validate() error {
	if !_regex_GetByNetworkRequest_NetworkId.MatchString(this.NetworkId) {
		return github_com_mwitkow_go_proto_validators.FieldError("NetworkId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.NetworkId))
	}
	if this.NetworkId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("NetworkId", fmt.Errorf(`value '%v' must not be an empty string`, this.NetworkId))
	}
	return nil
}
func (this *GetByNetworkResponse) Validate() error {
	for _, item := range this.Subscribers {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Subscribers", err)
			}
		}
	}
	return nil
}
func (this *DeleteSubscriberResponse) Validate() error {
	return nil
}

var _regex_GetSubscriberRequest_SubscriberId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *GetSubscriberRequest) Validate() error {
	if !_regex_GetSubscriberRequest_SubscriberId.MatchString(this.SubscriberId) {
		return github_com_mwitkow_go_proto_validators.FieldError("SubscriberId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.SubscriberId))
	}
	if this.SubscriberId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("SubscriberId", fmt.Errorf(`value '%v' must not be an empty string`, this.SubscriberId))
	}
	return nil
}
func (this *GetSubscriberResponse) Validate() error {
	if this.Subscriber != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Subscriber); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Subscriber", err)
		}
	}
	return nil
}

var _regex_GetSubscriberByEmailRequest_Email = regexp.MustCompile(`^$|^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func (this *GetSubscriberByEmailRequest) Validate() error {
	if !_regex_GetSubscriberByEmailRequest_Email.MatchString(this.Email) {
		return github_com_mwitkow_go_proto_validators.FieldError("Email", fmt.Errorf(`must be an email format`))
	}
	return nil
}
func (this *GetSubscriberByEmailResponse) Validate() error {
	if this.Subscriber != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Subscriber); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Subscriber", err)
		}
	}
	return nil
}

var _regex_AddSubscriberRequest_Email = regexp.MustCompile(`^$|^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
var _regex_AddSubscriberRequest_PhoneNumber = regexp.MustCompile(`^$|^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

func (this *AddSubscriberRequest) Validate() error {
	if this.Name == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Name", fmt.Errorf(`value '%v' must not be an empty string`, this.Name))
	}
	if !(len(this.Name) > 1) {
		return github_com_mwitkow_go_proto_validators.FieldError("Name", fmt.Errorf(`value '%v' must have a length greater than '1'`, this.Name))
	}
	if !_regex_AddSubscriberRequest_Email.MatchString(this.Email) {
		return github_com_mwitkow_go_proto_validators.FieldError("Email", fmt.Errorf(`must be an email format`))
	}
	if this.Email == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Email", fmt.Errorf(`must be an email format`))
	}
	if !_regex_AddSubscriberRequest_PhoneNumber.MatchString(this.PhoneNumber) {
		return github_com_mwitkow_go_proto_validators.FieldError("PhoneNumber", fmt.Errorf(`must be a phone number format`))
	}
	return nil
}

var _regex_UpdateSubscriberRequest_SubscriberId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)
var _regex_UpdateSubscriberRequest_PhoneNumber = regexp.MustCompile(`^$|^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

func (this *UpdateSubscriberRequest) Validate() error {
	if !_regex_UpdateSubscriberRequest_SubscriberId.MatchString(this.SubscriberId) {
		return github_com_mwitkow_go_proto_validators.FieldError("SubscriberId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.SubscriberId))
	}
	if this.SubscriberId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("SubscriberId", fmt.Errorf(`value '%v' must not be an empty string`, this.SubscriberId))
	}
	if !_regex_UpdateSubscriberRequest_PhoneNumber.MatchString(this.PhoneNumber) {
		return github_com_mwitkow_go_proto_validators.FieldError("PhoneNumber", fmt.Errorf(`must be a phone number format`))
	}
	if this.Name == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Name", fmt.Errorf(`value '%v' must not be an empty string`, this.Name))
	}
	if !(len(this.Name) > 1) {
		return github_com_mwitkow_go_proto_validators.FieldError("Name", fmt.Errorf(`value '%v' must have a length greater than '1'`, this.Name))
	}
	return nil
}
func (this *UpdateSubscriberResponse) Validate() error {
	return nil
}
func (this *AddSubscriberResponse) Validate() error {
	if this.Subscriber != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Subscriber); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Subscriber", err)
		}
	}
	return nil
}
