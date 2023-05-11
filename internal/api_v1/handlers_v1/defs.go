package handlers_v1

import "time"

const (
	TokenLength = 256

	UUIDLength = 36

	MinNicknameLength = 5
	MaxNicknameLength = 34

	MinEmailLength = 6
	MaxEmailLength = 64

	MinLoginLength = 6
	MaxLoginLength = 64

	MinPasswordLength = 6
	MaxPasswordLength = 64

	MinCouponLength = 4
	MaxCouponLength = 64

	MinTextLength = 3
	MaxTextLength = 64

	TempRegistrationExpiration = 10 * time.Minute
)
