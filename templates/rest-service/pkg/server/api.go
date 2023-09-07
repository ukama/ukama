package server

type FooGetRequest struct {
	Name string `path:"name" validate:"required"`
}

type FooGetResponse struct {
	Result string
}
