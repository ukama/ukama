package rest

// org group
type GetOrgsRequest struct {
	UserUuid string `example:"{{UserUUID}}" form:"user_uuid" json:"user_uuid" query:"user_uuid" binding:"required" validate:"required"`
}

type GetOrgRequest struct {
	OrgName string `example:"milky-way" path:"org" validate:"required"`
}

type GetByEmailRequest struct {
	Email string `example:" {{Email}}" path:"email" validate:"required"`
}

type AddOrgRequest struct {
	OrgName     string `example:"milky-way" json:"name" validate:"required"`
	Owner       string `example:"{{UserUUID}}" json:"owner_uuid" validate:"required"`
	Certificate string `example:"test_cert" json:"certificate"`
}

type UserOrgRequest struct {
	OrgId  string `example:"{{OrgId}}" path:"org_id"  validate:"required"`
	UserId string `example:"{{UserId}}" path:"user_id" validate:"required"`
}

// Users group

type GetUserRequest struct {
	UserId string `example:"{{UserID}}" path:"user_id" validate:"required"`
}

type GetUserByAuthIdRequest struct {
	AuthId string `example:"{{AuthId}}" path:"auth_id" validate:"required"`
}

type AddUserRequest struct {
	Name   string `example:"John" json:"name,omitempty" validate:"required"`
	Email  string `example:"john@example.com" json:"email" validate:"required"`
	Phone  string `example:"4151231234" json:"phone,omitempty"`
	AuthId string `example:"{{AuthId}}" json:"auth_id" validate:"required"`
}
