package handlers_v1

import (
	"encoding/json"
	"net/http"
	"strings"
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
	var input struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decodeErr := json.NewDecoder(r.Body).Decode(&input)
	if decodeErr != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	// Block 1 - data validation
	nickname := strings.TrimSpace(input.Nickname)
	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := strings.TrimSpace(input.Password)

	if err := tl.Validate(nickname, tl.IsNotBlank(), tl.IsMinMaxLen(MinNicknameLength, MaxNicknameLength), tl.IsNotContainsSpace(), tl.IsNickname()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Nickname: "+err.Error())
		return
	}
	if err := tl.Validate(email, tl.IsNotBlank(), tl.IsMinMaxLen(MinEmailLength, MaxEmailLength), tl.IsNotContainsSpace(), tl.IsEmail()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Email: "+err.Error())
		return
	}
	if err := tl.Validate(password, tl.IsNotBlank(), tl.IsMinMaxLen(MinPasswordLength, MaxPasswordLength), tl.IsNotContainsSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Password: "+err.Error())
		return
	}

	emailDomainExists, err := tl.CheckEmailDomainExistence(email)
	if !emailDomainExists {
		api_v1.RespondWithConflict(w, "Email: the email domain is not exist")
		return
	}
	if err != nil {
		rs.App.Logger.NewWarn("Error in checked the email domain: ", err)
	}

	// Block 2 - check for an exists nickname and email
	nicknameExist, emailExist, err := storage.CheckUserExists(r.Context(), rs.App.Postgres, nickname, email)
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

	if err := storage.CreateTempRegistration(r.Context(), rs.App.Redis, nickname, email, password, confirmationUrlToken, TempRegistrationExpiration); err != nil {
		rs.App.Logger.NewWarn("error in inserted registration temp record", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 4 - generate url and send url on email
	url, err := tl.UrlSetParam(rs.App.Config.App.Service.Url.App+"/confirm-signup", "token", confirmationUrlToken)
	if err != nil {
		rs.App.Logger.NewWarn("error in url set param", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	if err = rs.App.Mailer.SendEmailConfirmation(nickname, email, url); err != nil {
		rs.App.Logger.NewWarn("error in sent email confirmation", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 5 - send the result
	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AuthSignupWithToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Token string `json:"token"`
	}
	decodeErr := json.NewDecoder(r.Body).Decode(&input)
	if decodeErr != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	// Block 1 - data validation
	token := strings.TrimSpace(input.Token)

	if err := tl.Validate(token, tl.IsNotBlank(), tl.IsLen(TokenLength), tl.IsNotContainsSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Token: "+err.Error())
		return
	}

	// Block 2 - get user data and check on exist user
	nickname, email, password, err := storage.GetTempRegistration(r.Context(), rs.App.Redis, token)
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

	userUuid, err := storage.CreateUser(r.Context(), rs.App.Postgres, rs.App.Redis, nickname, email, base64PasswordHash, base64Salt, token)
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
		AvatarUrl:          rs.App.Config.App.Service.Url.Api + storage.ResourcesProfileImagePath + tl.UuidToStringNoDashes(userUuid),
	}

	api_v1.RespondWithCreated(w, response)
}

func (rs *Resolver) AuthLogin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nickname *string `json:"nickname"`
		Email    *string `json:"email"`
		Password string  `json:"password"`
	}
	decodeErr := json.NewDecoder(r.Body).Decode(&input)
	if decodeErr != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	// Block 1 - data validation
	var nickname, email string
	if input.Nickname != nil {
		nickname = strings.TrimSpace(strings.ToLower(*input.Nickname))
	}
	if input.Email != nil {
		email = strings.TrimSpace(strings.ToLower(*input.Email))
	}
	password := strings.TrimSpace(input.Password)

	if err := tl.Validate(password, tl.IsNotBlank(), tl.IsMinMaxLen(MinPasswordLength, MaxPasswordLength), tl.IsNotContainsSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Password: "+err.Error())
		return
	}

	if nickname == "" && email == "" {
		api_v1.RespondWithUnprocessableEntity(w, "Nickname and Email: the values is empty")
		return
	} else if nickname != "" {
		if err := tl.Validate(nickname, tl.IsNotBlank(), tl.IsMinMaxLen(MinNicknameLength, MaxNicknameLength), tl.IsNotContainsSpace(), tl.IsNickname()); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Nickname: "+err.Error())
			return
		}
	} else {
		if err := tl.Validate(email, tl.IsNotBlank(), tl.IsMinMaxLen(MinEmailLength, MaxEmailLength), tl.IsNotContainsSpace(), tl.IsEmail()); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Email: "+err.Error())
			return
		}

		emailDomainExists, err := tl.CheckEmailDomainExistence(email)
		if !emailDomainExists {
			api_v1.RespondWithConflict(w, "Email: the email domain is not exist")
			return
		} else if err != nil {
			rs.App.Logger.NewWarn("Error in checked the email domain: ", err)
		}
	}

	// Block 2 - check for an exists nickname and email
	nicknameExist, emailExist, err := storage.CheckUserExists(r.Context(), rs.App.Postgres, nickname, email)
	if err != nil {
		rs.App.Logger.NewWarn("error in checked the user existence", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}
	if !nicknameExist && !emailExist {
		api_v1.RedRespond(w, http.StatusNotFound, "Not found", "There is no user with this nickname and email address")
		return
	} else {
		userUuid, scannedNickname, scannedEmail, base64PasswordHash, base64Salt, err := storage.GetUserData(r.Context(), rs.App.Postgres, nickname, email)
		if err != nil {
			rs.App.Logger.NewWarn("error in get user password and salt", err)
			api_v1.RespondWithInternalServerError(w)
			return
		}

		// Get account state and check on exists
		state, err := storage.GetStateAccount(r.Context(), rs.App.Postgres, userUuid, storage.AccountRoleUser)
		if state == "" {
			api_v1.RedRespond(w, http.StatusUnauthorized, "Unauthorized", "The account was not found in the list of users")
			return
		} else if err != nil {
			api_v1.RespondWithInternalServerError(w)
			rs.App.Logger.NewWarn("Error in founding account in the list", err)
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

		result, err := auth.CompareHashPasswords(password, base64PasswordHash, base64Salt)
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
			AvatarUrl:          rs.App.Config.App.Service.Url.Api + storage.ResourcesProfileImagePath + userUuid,
		}

		api_v1.RespondWithCreated(w, response)
	}
}

func (rs *Resolver) AuthLogout(w http.ResponseWriter, r *http.Request) {
	token, data, err := api_v1.ContextGetAuthenticated(r)
	if err != nil {
		rs.App.Logger.NewWarn("error in took jwt data", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	ttl := data.ExpiresAt.Sub(time.Now())
	if err := storage.CreateBlockedToken(r.Context(), rs.App.Redis, token, ttl); err == storage.QueryExists {
		api_v1.RedRespond(w, http.StatusUnauthorized, "Unauthorized", "Token has already been deactivated")
		return
	} else if err != nil {
		rs.App.Logger.NewWarn("error in took jwt data", err)
		api_v1.RespondWithInternalServerError(w)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AuthLoginWithToken(w http.ResponseWriter, r *http.Request)           {}
func (rs *Resolver) AuthRecoverPassword(w http.ResponseWriter, r *http.Request)          {}
func (rs *Resolver) AuthRecoverPasswordWithToken(w http.ResponseWriter, r *http.Request) {}
