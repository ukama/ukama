package rest

type GetUserInfo struct {
	Id    string `example:"1" json:"id" validation:"required"`
	Name  string `example:"John Doe" json:"name" validation:"required"`
	Email string `example:"john@example.com" json:"email" validation:"required"`
}

type Authenticate struct {
	IsValidSession bool `example:"1" json:"is_valid_session" validation:"required"`
}
