// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: distributor.proto

package gen

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *NotificationsRequest) Validate() error {
	return nil
}
func (this *NotificationsResponse) Validate() error {
	for _, item := range this.Notifications {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Notifications", err)
			}
		}
	}
	return nil
}
func (this *NotificationStreamRequest) Validate() error {
	return nil
}
func (this *Notification) Validate() error {
	return nil
}
