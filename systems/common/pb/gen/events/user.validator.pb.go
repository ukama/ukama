// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: events/user.proto

package events

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/mwitkow/go-proto-validators"
	regexp "regexp"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

var _regex_EventUserCreate_UserId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *EventUserCreate) Validate() error {
	if !_regex_EventUserCreate_UserId.MatchString(this.UserId) {
		return github_com_mwitkow_go_proto_validators.FieldError("UserId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.UserId))
	}
	if this.UserId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("UserId", fmt.Errorf(`value '%v' must not be an empty string`, this.UserId))
	}
	return nil
}

var _regex_EventUserUpdate_UserId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *EventUserUpdate) Validate() error {
	if !_regex_EventUserUpdate_UserId.MatchString(this.UserId) {
		return github_com_mwitkow_go_proto_validators.FieldError("UserId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.UserId))
	}
	if this.UserId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("UserId", fmt.Errorf(`value '%v' must not be an empty string`, this.UserId))
	}
	return nil
}

var _regex_EventUserDeactivate_UserId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *EventUserDeactivate) Validate() error {
	if !_regex_EventUserDeactivate_UserId.MatchString(this.UserId) {
		return github_com_mwitkow_go_proto_validators.FieldError("UserId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.UserId))
	}
	if this.UserId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("UserId", fmt.Errorf(`value '%v' must not be an empty string`, this.UserId))
	}
	return nil
}

var _regex_EventUserDelete_UserId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *EventUserDelete) Validate() error {
	if !_regex_EventUserDelete_UserId.MatchString(this.UserId) {
		return github_com_mwitkow_go_proto_validators.FieldError("UserId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.UserId))
	}
	if this.UserId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("UserId", fmt.Errorf(`value '%v' must not be an empty string`, this.UserId))
	}
	return nil
}
