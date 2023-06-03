package storage

import (
	"errors"
	"strings"
)

// Names style:
// For creating a record: Create<Type>
// For checking the existence of a record: Check<Type>
// For getting a record: Get<Type>
// For updating a record: Update<Type>
// For deleting a record: Delete<Type>

// Accounts
const (
	AccountStateActive  = "active"
	AccountStateBlocked = "blocked"
	AccountStateDeleted = "deleted"

	AccountRegistrationMethodWebApplication  = "web application"
	AccountRegistrationMethodFromAdminPanel  = "from admin panel"
	AccountRegistrationMethodTelegramAccount = "telegram account"
	AccountRegistrationMethodGoogleAccount   = "google account"

	AccountRoleUser  = "user"
	AccountRoleAdmin = "admin"
)

// Products
const (
	ProductStateActive                  = "active"
	ProductStateUnavailableWithPrice    = "unavailable with price"
	ProductStateUnavailableWithoutPrice = "unavailable without price"
	ProductStateInvisible               = "invisible"
	ProductStateDeleted                 = "deleted"
)

// System
const (
	ResourcesProfileImagePath = "/api/v1/profile/image/"
	ResourcesProductImagePath = "/api/v1/resources/product_image/"
	ResourcesSvgFilePath      = "/api/v1/resources/svg/"
)

var (
	FailedUpdate = errors.New("failed to update data")
	FailedInsert = errors.New("failed to insert data")
	FailedDelete = errors.New("failed to delete data")

	NoResults   = errors.New("no results")
	QueryExists = errors.New("value already exists")
)

func GetProfileImageUrl(apiUrl, file string) string {
	return strings.ToLower(apiUrl + ResourcesProfileImagePath + strings.ReplaceAll(file, " ", "-"))
}

func GetProductImageUrl(apiUrl, file string) string {
	return strings.ToLower(apiUrl + ResourcesProductImagePath + strings.ReplaceAll(file, " ", "-"))
}

func GetSvgFileUrl(apiUrl, file string) string {
	return strings.ToLower(apiUrl + ResourcesSvgFilePath + strings.ReplaceAll(file, " ", "-"))
}

// Errors
const (
	PgNoRows     = "not found value(s)"
	PgNoUpdated  = "the data not updated"
	PgForeignKey = "the value is using now"
	PgNoUnique   = "the value not unique"
)
