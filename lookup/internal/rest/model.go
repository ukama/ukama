package rest

import (
	uuid "github.com/satori/go.uuid"
)

type ErrorMessage struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type DeviceMappingRequest struct {
	Org string `json:"org" binding:"required"`
}

type GetDeviceResponse struct {
	Uuid        uuid.UUID `json:"uuid"`
	OrgName     string    `json:"orgName"`
	Certificate string    `json:"certificate" binding:"base64"`
	Ip          string    `json:"ip" validate:"ip"`
}

type AddOrgRequest struct {
	Certificate string `json:"certificate" binding:"required,base64"`
	Ip          string `json:"ip" binding:"required,ip"`
}
