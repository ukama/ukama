// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: member.proto

package gen

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/mwitkow/go-proto-validators"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
	_ "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	regexp "regexp"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

var _regex_GetMemberByUserIdRequest_MemberId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *GetMemberByUserIdRequest) Validate() error {
	if !_regex_GetMemberByUserIdRequest_MemberId.MatchString(this.MemberId) {
		return github_com_mwitkow_go_proto_validators.FieldError("MemberId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.MemberId))
	}
	if this.MemberId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("MemberId", fmt.Errorf(`value '%v' must not be an empty string`, this.MemberId))
	}
	return nil
}
func (this *GetMemberByUserIdResponse) Validate() error {
	if this.Member != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Member); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Member", err)
		}
	}
	return nil
}

var _regex_AddMemberRequest_UserUuid = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *AddMemberRequest) Validate() error {
	if !_regex_AddMemberRequest_UserUuid.MatchString(this.UserUuid) {
		return github_com_mwitkow_go_proto_validators.FieldError("UserUuid", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.UserUuid))
	}
	if this.UserUuid == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("UserUuid", fmt.Errorf(`value '%v' must not be an empty string`, this.UserUuid))
	}
	return nil
}

var _regex_MemberRequest_MemberId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *MemberRequest) Validate() error {
	if !_regex_MemberRequest_MemberId.MatchString(this.MemberId) {
		return github_com_mwitkow_go_proto_validators.FieldError("MemberId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.MemberId))
	}
	if this.MemberId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("MemberId", fmt.Errorf(`value '%v' must not be an empty string`, this.MemberId))
	}
	return nil
}
func (this *MemberResponse) Validate() error {
	if this.Member != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Member); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Member", err)
		}
	}
	return nil
}
func (this *GetMembersRequest) Validate() error {
	return nil
}
func (this *GetMembersResponse) Validate() error {
	for _, item := range this.Members {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Members", err)
			}
		}
	}
	return nil
}

var _regex_UpdateMemberRequest_MemberId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *UpdateMemberRequest) Validate() error {
	if !_regex_UpdateMemberRequest_MemberId.MatchString(this.MemberId) {
		return github_com_mwitkow_go_proto_validators.FieldError("MemberId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.MemberId))
	}
	if this.MemberId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("MemberId", fmt.Errorf(`value '%v' must not be an empty string`, this.MemberId))
	}
	return nil
}

var _regex_Member_MemberId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)
var _regex_Member_UserId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *Member) Validate() error {
	if !_regex_Member_MemberId.MatchString(this.MemberId) {
		return github_com_mwitkow_go_proto_validators.FieldError("MemberId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.MemberId))
	}
	if this.MemberId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("MemberId", fmt.Errorf(`value '%v' must not be an empty string`, this.MemberId))
	}
	if !_regex_Member_UserId.MatchString(this.UserId) {
		return github_com_mwitkow_go_proto_validators.FieldError("UserId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.UserId))
	}
	if this.UserId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("UserId", fmt.Errorf(`value '%v' must not be an empty string`, this.UserId))
	}
	if this.CreatedAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.CreatedAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("CreatedAt", err)
		}
	}
	return nil
}
