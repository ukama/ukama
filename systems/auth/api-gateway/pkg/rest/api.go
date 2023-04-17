package rest

type GetUserInfo struct {
	Id         string `example:"1" json:"id" validation:"required"`
	Name       string `example:"John Doe" json:"name" validation:"required"`
	Email      string `example:"john@example.com" json:"email" validation:"required"`
	Role       string `example:"admin" json:"role" validation:"required"`
	FirstVisit bool   `example:"true" json:"first_visit" validation:"required"`
}

type OptionalReqHeader struct {
	XSessionToken string `example:"" json:"x-session-token" header:"x-session-token" validation:"required" required:"false"`
}

type LoginReq struct {
	Email    string `example:"john@example.com" json:"email" validation:"required"`
	Password string `example:"Password" json:"password" validation:"required"`
}
type LoginRes struct {
	Token string `example:"Token" json:"token" validation:"required"`
}

type GetSessionReq struct {
	XSessionToken string `example:"token" json:"x-session-token" header:"x-session-token" validation:"required"`
}
