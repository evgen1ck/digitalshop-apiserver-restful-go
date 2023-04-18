package handlers

import "time"

const (
	TokenLength = 256

	MinNicknameLength = 5
	MaxNicknameLength = 34

	MinEmailLength = 6
	MaxEmailLength = 64

	MinPasswordLength = 6
	MaxPasswordLength = 64

	TempRegistrationExpiration = 10 * time.Minute
)
