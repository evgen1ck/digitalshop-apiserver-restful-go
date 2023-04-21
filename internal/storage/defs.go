package storage

import (
	"errors"
)

const (
	AccountRoleUser  = "user"
	AccountRoleAdmin = "admin"

	AccountStateActive  = "active"
	AccountStateBlocked = "blocked"
	AccountStateDeleted = "deleted"

	AccountRegistrationMethodWebApplication  = "web application"
	AccountRegistrationMethodTelegramAccount = "telegram account"
	AccountRegistrationMethodGoogleAccount   = "google account"

	ResourcesProfileImagePath = "/api/v1/profile/image/"
	ResourcesProductImagePath = "/api/v1/resources/product_image/"
	ResourcesSvgFilePath      = "/api/v1/resources/svg_file/"
)

var (
	FailedUpdate = errors.New("failed to update data")
	FailedInsert = errors.New("failed to insert data")
	FailedDelete = errors.New("failed to delete data")

	NoResults   = errors.New("no results")
	QueryExists = errors.New("value already exists")
)
