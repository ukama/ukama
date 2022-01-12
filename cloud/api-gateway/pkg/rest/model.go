package rest

type UserRequest struct {
	Org       string `path:"org" validate:"required"`
	Imsi      string `json:"imsi" validate:"required"`
	FirstName string `json:"firstName,omitempty" validate:"required"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
}
