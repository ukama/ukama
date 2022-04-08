package server

type ReqAddUpdateNode struct {
	Nodeid     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespAddUpdateNode struct {
}

type ReqUpdateNodeStatus struct {
	Nodeid     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespUpdateNodeStatus struct {
}

type ReqDeleteNode struct {
	Nodeid     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespDeleteNode struct {
}

type ReqAddUpdateModule struct {
	Nodeid     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type RespAddUpdateModule struct {
}

type ReqUpdateBootStrapCerts struct {
	Nodeid     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type ReqUpdateUserConfig struct {
	Nodeid     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type ReqUpdateFactoryConfig struct {
	Nodeid     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type ReqUpdateUserCalibration struct {
	Nodeid     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}

type ReqUpdateFactoryCalibration struct {
	Nodeid     string `query:"nodeid" validate:"required"`
	LookingFor string `query:"looking_for" validate:"required"`
}
