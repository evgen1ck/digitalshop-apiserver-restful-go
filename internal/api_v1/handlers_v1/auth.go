package handlers_v1

import (
	"encoding/json"
	"net/http"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/auth"
	"test-server-go/internal/storage"
	tl "test-server-go/internal/tools"
	"time"
)

type authResponse struct {
	Token              string `json:"token"`
	Role               string `json:"role"`
	Uuid               string `json:"uuid"`
	Nickname           string `json:"nickname"`
	Email              string `json:"email"`
	RegistrationMethod string `json:"registration_method"`
	AvatarUrl          string `json:"avatar_url"`
}

func (rs *Resolver) AuthSignup(w http.ResponseWriter, r *http.Request) {
	// Block 0 - decode data
	var data struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decodeErr := json.NewDecoder(r.Body).Decode(&data)
	if decodeErr != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	// Block 1 - data validation
	if err := tl.Validate(data.Nickname, tl.IsNotBlank(), tl.IsMinMaxLen(MinNicknameLength, MaxNicknameLength), tl.IsNotContainsSpace(), tl.IsNickname(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Nickname: "+err.Error())
		return
	}
	if err := tl.Validate(data.Email, tl.IsNotBlank(), tl.IsMinMaxLen(MinEmailLength, MaxEmailLength), tl.IsNotContainsSpace(), tl.IsEmail(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Email: "+err.Error())
		return
	}
	if err := tl.Validate(data.Password, tl.IsNotBlank(), tl.IsMinMaxLen(MinPasswordLength, MaxPasswordLength), tl.IsNotContainsSpace(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Password: "+err.Error())
		return
	}

	emailDomainExists, err := tl.CheckEmailDomainExistence(data.Email)
	if err != nil {
		rs.App.Logger.NewWarn("Error in checked the email domain: ", err)
	} else if !emailDomainExists {
		api_v1.RespondWithConflict(w, "Email: the email domain is not exist")
		return
	}

	// Block 2 - check for an exists nickname and email
	nicknameExist, emailExist, err := storage.CheckUser(r.Context(), rs.App.Postgres, data.Nickname, data.Email)
	if err != nil {
		rs.App.Logger.NewWarn("error in checked the user existence", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}
	if nicknameExist {
		api_v1.RespondWithConflict(w, "Nickname: this nickname is already in use")
		return
	}
	if emailExist {
		api_v1.RespondWithConflict(w, "Email: this email is already in use")
		return
	}

	// Block 3 - generate token and insert a temporary account record
	confirmationUrlToken, err := tl.GenerateURLToken(TokenLength)
	if err != nil {
		rs.App.Logger.NewWarn("error in generated url token", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	if err = storage.CreateTempRegistration(r.Context(), rs.App.Redis, data.Nickname, data.Email, data.Password, confirmationUrlToken, TempRegistrationExpiration); err != nil {
		rs.App.Logger.NewWarn("error in inserted registration temp record", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 4 - generate url and send url on email
	url, err := tl.UrlSetParam(rs.App.Config.App.Service.Url.Client+"/confirm-signup", "token", confirmationUrlToken)
	if err != nil {
		rs.App.Logger.NewWarn("error in url set param", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	if err = rs.App.Mailer.SendEmailConfirmation(data.Nickname, data.Email, url); err != nil {
		rs.App.Logger.NewWarn("error in sent email confirmation", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 5 - send the result
	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AuthSignupWithToken(w http.ResponseWriter, r *http.Request) {
	// Block 0 - decode data
	var data struct {
		Token string `json:"token"`
	}
	decodeErr := json.NewDecoder(r.Body).Decode(&data)
	if decodeErr != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	// Block 1 - data validation
	if err := tl.Validate(data.Token, tl.IsNotBlank(), tl.IsLen(TokenLength), tl.IsNotContainsSpace(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Token: "+err.Error())
		return
	}

	// Block 2 - get user data and check on exist user
	nickname, email, password, err := storage.GetTempRegistration(r.Context(), rs.App.Redis, data.Token)
	if password == "" {
		api_v1.RespondWithConflict(w, "User not found")
		return
	} else if err != nil {
		rs.App.Logger.NewWarn("error in checked registration temp record", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 3 - hash password and add a user
	base64PasswordHash, base64Salt, err := auth.HashPassword(password, "")
	if err != nil {
		rs.App.Logger.NewWarn("error in generated hash password", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	userUuid, err := storage.CreateUser(r.Context(), rs.App.Postgres, rs.App.Redis, nickname, email, base64PasswordHash, base64Salt, data.Token)
	if err != nil {
		rs.App.Logger.NewWarn("error in registration user", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 4 - generate JWT
	jwtToken, err := auth.GenerateJwt(userUuid, rs.App.Config.App.Jwt)
	if err != nil {
		rs.App.Logger.NewWarn("error in generated jwt", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 5 - send the result
	response := authResponse{
		Token:              jwtToken,
		Role:               storage.AccountRoleUser,
		Uuid:               userUuid,
		Nickname:           nickname,
		Email:              email,
		RegistrationMethod: storage.AccountRegistrationMethodWebApplication,
		AvatarUrl:          rs.App.Config.App.Service.Url.Server + storage.ResourcesProfileImagePath + userUuid,
	}

	api_v1.RespondWithCreated(w, response)
}

func (rs *Resolver) AuthLogin(w http.ResponseWriter, r *http.Request) {
	// Block 0 - decode data
	var data struct {
		Nickname *string `json:"nickname"`
		Email    *string `json:"email"`
		Password string  `json:"password"`
	}
	decodeErr := json.NewDecoder(r.Body).Decode(&data)
	if decodeErr != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	// Block 1 - data validation
	if err := tl.Validate(data.Password, tl.IsNotBlank(), tl.IsMinMaxLen(MinPasswordLength, MaxPasswordLength), tl.IsNotContainsSpace(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Password: "+err.Error())
		rs.App.Logger.NewWarn("Password: ", err)
		return
	}
	var nickname, email string
	if data.Nickname != nil && *data.Nickname != "" {
		nickname = *data.Nickname
		if err := tl.Validate(nickname, tl.IsNotBlank(), tl.IsMinMaxLen(MinNicknameLength, MaxNicknameLength), tl.IsNotContainsSpace(), tl.IsNickname(), tl.IsTrimmedSpace()); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Nickname: "+err.Error())
			rs.App.Logger.NewWarn("Nickname: ", err)
			return
		}
	} else if data.Email != nil && *data.Email != "" {
		email = *data.Email
		if err := tl.Validate(email, tl.IsNotBlank(), tl.IsMinMaxLen(MinEmailLength, MaxEmailLength), tl.IsNotContainsSpace(), tl.IsEmail(), tl.IsTrimmedSpace()); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Email: "+err.Error())
			rs.App.Logger.NewWarn("Email: ", err)
			return
		}
		emailDomainExists, err := tl.CheckEmailDomainExistence(email)
		if err != nil {
			rs.App.Logger.NewWarn("Error in checked the email domain: ", err)
		} else if !emailDomainExists {
			api_v1.RespondWithConflict(w, "Email: the email domain is not exist")
			return
		}
	} else {
		api_v1.RespondWithUnprocessableEntity(w, "Nickname and Email: the values are empty")
		rs.App.Logger.NewWarn("Nickname and Email: the values are empty ", nil)
		return
	}

	// Block 2 - check for an exists nickname and email
	nicknameExist, emailExist, err := storage.CheckUser(r.Context(), rs.App.Postgres, nickname, email)
	if err != nil {
		rs.App.Logger.NewWarn("error in checked the user existence", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}
	if !nicknameExist && !emailExist {
		if nickname == "" {
			api_v1.RedRespond(w, http.StatusNotFound, "Not found", "User with this email was not found")
			return
		} else {
			api_v1.RedRespond(w, http.StatusNotFound, "Not found", "User with this nickname was not found")
			return
		}
	}

	userUuid, scannedNickname, scannedEmail, base64PasswordHash, base64Salt, err := storage.GetUserData(r.Context(), rs.App.Postgres, nickname, email)
	if err != nil {
		rs.App.Logger.NewWarn("error in get user password and salt", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Get account state and check on exists
	state, err := storage.GetStateAccount(r.Context(), rs.App.Postgres, userUuid, storage.AccountRoleUser)
	if err != nil {
		api_v1.RespondWithInternalServerError(w)
		rs.App.Logger.NewWarn("Error in founding account in the list", err)
		return
	} else if state == "" {
		api_v1.RedRespond(w, http.StatusUnauthorized, "Unauthorized", "The account was not found in the list of users")
		return
	}

	// Check account on state (blocked, deleted...)
	switch state {
	case storage.AccountStateBlocked:
		api_v1.RedRespond(w, http.StatusForbidden, "Forbidden", "This account has been blocked")
		return
	case storage.AccountStateDeleted:
		api_v1.RedRespond(w, http.StatusForbidden, "Forbidden", "This account has been deleted")
		return
	}

	result, err := auth.CompareHashPasswords(data.Password, base64PasswordHash, base64Salt)
	if err != nil {
		rs.App.Logger.NewWarn("error in compare hash passwords", err)
		api_v1.RespondWithInternalServerError(w)
		return
	} else if result == false {
		api_v1.RedRespond(w, http.StatusUnauthorized, "Unauthorized", "Invalid password")
		return
	}

	// Block 4 - generate JWT
	jwtToken, err := auth.GenerateJwt(userUuid, rs.App.Config.App.Jwt)
	if err != nil {
		rs.App.Logger.NewWarn("error in generated jwt", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 5 - send the result
	response := authResponse{
		Token:              jwtToken,
		Role:               storage.AccountRoleUser,
		Uuid:               userUuid,
		Nickname:           scannedNickname,
		Email:              scannedEmail,
		RegistrationMethod: storage.AccountRegistrationMethodWebApplication,
		AvatarUrl:          rs.App.Config.App.Service.Url.Server + storage.ResourcesProfileImagePath + userUuid,
	}

	api_v1.RespondWithCreated(w, response)
}

func (rs *Resolver) AuthLogout(w http.ResponseWriter, r *http.Request) {
	// Block 0 - decode data
	token, data, err := api_v1.ContextGetAuthenticated(r)
	if err != nil {
		rs.App.Logger.NewWarn("error in took jwt data", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 1 - add token in stop-list
	ttl := data.ExpiresAt.Sub(time.Now())
	if err = storage.CreateBlockedToken(r.Context(), rs.App.Redis, token, ttl); err == storage.QueryExists {
		api_v1.RedRespond(w, http.StatusUnauthorized, "Unauthorized", "Token has already been deactivated")
		return
	} else if err != nil {
		rs.App.Logger.NewWarn("error in took jwt data", err)
		api_v1.RespondWithInternalServerError(w)
	}

	// Block 2 - send the result
	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AuthLoginWithToken(w http.ResponseWriter, r *http.Request)           {}
func (rs *Resolver) AuthRecoverPassword(w http.ResponseWriter, r *http.Request)          {}
func (rs *Resolver) AuthRecoverPasswordWithToken(w http.ResponseWriter, r *http.Request) {}
