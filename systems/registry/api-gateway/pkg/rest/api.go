package rest

// org group

type GetOrgRequest struct {
	OrgName string `path:"org" validate:"required"`
}

type AddOrgRequest struct {
	OrgName     string `json:"name" validate:"required"`
	Owner       string `json:"owner_uuid" validate:"required"`
	Certificate string `json:"certificate"`
}

type MemberRequest struct {
	OrgName  string `path:"org" validate:"required"`
	UserUUID string `json:"user_uuid" validate:"required"`
}

type GetMemberRequest struct {
	OrgName  string `path:"org" validate:"required"`
	UserUUID string `path:"user_uuid" validate:"required"`
}

// Users group

type GetUserRequest struct {
	UserUUID string `path:"user_uuid" validate:"required"`
}

type UserRequest struct {
	Org      string `path:"org" validate:"required"`
	SimToken string `json:"simToken"`
	Name     string `json:"name,omitempty" validate:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"`
}

type UpdateUserRequest struct {
	OrgName       string `path:"org" validate:"required"`
	UserId        string `path:"user" validate:"required"`
	Name          string `json:"name,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	IsDeactivated bool   `json:"isDeactivated,omitempty"`
}

type SetSimStatusRequest struct {
	OrgName string       `path:"org" validate:"required"`
	UserId  string       `path:"user" validate:"required"`
	Iccid   string       `path:"iccid" validate:"required"`
	Carrier *SimServices `json:"carrier,omitempty"`
	Ukama   *SimServices `json:"ukama,omitempty"`
}

type GetSimQrRequest struct {
	OrgName string `path:"org" validate:"required"`
	UserId  string `path:"user" validate:"required"`
	Iccid   string `path:"iccid" validate:"required"`
}

type SimServices struct {
	Voice *bool `json:"voice,omitempty"`
	Sms   *bool `json:"sms,omitempty"`
	Data  *bool `json:"data,omitempty"`
}

type DeleteUserRequest struct {
	OrgName string `path:"org" validate:"required"`
	UserId  string `path:"user" validate:"required"`
}
