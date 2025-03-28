// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mailer.proto

package gen

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	_ "github.com/mwitkow/go-proto-validators"
	regexp "regexp"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *GetEmailByIdRequest) Validate() error {
	return nil
}

var _regex_GetEmailByIdResponse_MailId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *GetEmailByIdResponse) Validate() error {
	if !_regex_GetEmailByIdResponse_MailId.MatchString(this.MailId) {
		return github_com_mwitkow_go_proto_validators.FieldError("MailId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.MailId))
	}
	if this.MailId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("MailId", fmt.Errorf(`value '%v' must not be an empty string`, this.MailId))
	}
	// Validation of proto3 map<> fields is unsupported.
	if this.CreatedAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.CreatedAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("CreatedAt", err)
		}
	}
	if this.UpdatedAt != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.UpdatedAt); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("UpdatedAt", err)
		}
	}
	return nil
}
func (this *SendEmailRequest) Validate() error {
	if this.TemplateName == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("TemplateName", fmt.Errorf(`value '%v' must not be an empty string`, this.TemplateName))
	}
	// Validation of proto3 map<> fields is unsupported.
	for _, item := range this.Attachments {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Attachments", err)
			}
		}
	}
	return nil
}
func (this *Attachment) Validate() error {
	return nil
}

var _regex_SendEmailResponse_MailId = regexp.MustCompile(`^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$`)

func (this *SendEmailResponse) Validate() error {
	if !_regex_SendEmailResponse_MailId.MatchString(this.MailId) {
		return github_com_mwitkow_go_proto_validators.FieldError("MailId", fmt.Errorf(`value '%v' must be a string conforming to regex "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[4][a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})?$"`, this.MailId))
	}
	if this.MailId == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("MailId", fmt.Errorf(`value '%v' must not be an empty string`, this.MailId))
	}
	return nil
}
