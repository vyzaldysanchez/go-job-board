package models

import "github.com/gbrlsnchs/jwt/v3"

type CustomPayload struct {
	jwt.Payload
	Email string `json:"email,omitempty"`
}
