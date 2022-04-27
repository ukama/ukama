package server

type ChunkRequest struct {
	Name    string `path:"name" validate:"required"`
	Version string `path:"version" validate:"required"`
	Store   string `json:"store" validate:"required"`
}

type ChunkResponse struct {
}

type CApp struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
