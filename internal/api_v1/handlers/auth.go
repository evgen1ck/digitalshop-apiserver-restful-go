package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/auth"
	"test-server-go/internal/storage"
	tl "test-server-go/internal/tools"
)

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

	if err := tl.Validate(nickname, tl.IsNotBlank(), tl.IsMinMaxLen(5, 32), tl.IsNotContainsSpace(), tl.IsNickname()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Nickname: "+err.Error())
		return
	}
	if err := tl.Validate(email, tl.IsNotBlank(), tl.IsMinMaxLen(6, 64), tl.IsNotContainsSpace(), tl.IsEmail()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Email: "+err.Error())
		return
	}
	if err := tl.Validate(password, tl.IsNotBlank(), tl.IsMinMaxLen(6, 64), tl.IsNotContainsSpace()); err != nil {
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

	// Block 2 - checking for an existing nickname and email
	nicknameExist, emailExist, err := storage.CheckUserExists(r.Context(), rs.App.Postgres, nickname, email)
	if err != nil {
		rs.App.Logger.NewError("error in checked the user existence", err)
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

	// Block 3 - generating token and inserting a temporary account record
	confirmationUrlToken, err := tl.GenerateURLToken(256)
	if err != nil {
		rs.App.Logger.NewError("error in generated url token", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	err = storage.CreateTempRegistration(r.Context(), rs.App.Postgres, nickname, email, password, confirmationUrlToken)
	if err != nil {
		rs.App.Logger.NewError("error in inserted registration temp record", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 4 - generating url and sending url on email
	url, err := tl.UrlSetParam(rs.App.Config.App.Service.Url.App+"/confirm-signup", "token", confirmationUrlToken)
	if err != nil {
		rs.App.Logger.NewError("error in url set param", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	err = rs.App.Mailer.SendEmailConfirmation(nickname, email, url)
	if err != nil {
		rs.App.Logger.NewError("error in sent email confirmation", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 5 - sending the result
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

	if err := tl.Validate(token, tl.IsNotBlank(), tl.IsLen(256), tl.IsNotContainsSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Token: "+err.Error())
		return
	}

	// Block 2 - get user data and checking on exist user
	nickname, email, password, err := storage.GetTempRegistration(r.Context(), rs.App.Postgres, token)
	if password == "" {
		api_v1.RespondWithConflict(w, "User not found")
		return
	}
	if err != nil {
		rs.App.Logger.NewError("error in checked registration temp record", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 3 - hashing password and adding a user
	base64PasswordHash, base64Salt, err := auth.HashPassword(password, "")
	if err != nil {
		rs.App.Logger.NewError("error in generated hash password", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	userUuid, err := storage.CreateUser(r.Context(), rs.App.Postgres, nickname, email, base64PasswordHash, base64Salt)
	if err != nil {
		rs.App.Logger.NewError("error in registration user", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 4 - generating JWT
	jwt, err := auth.GenerateJwt(userUuid, rs.App.Config.App.Jwt)
	if err != nil {
		rs.App.Logger.NewError("error in generated jwt", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 5 - sending the result
	response := struct {
		Token    string `json:"token"`
		Uuid     string `json:"uuid"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
	}{
		Token:    jwt,
		Uuid:     userUuid,
		Nickname: nickname,
		Email:    email,
	}
	api_v1.RespondWithCreated(w, response)
}

func (rs *Resolver) AuthLogin(w http.ResponseWriter, r *http.Request)                    {}
func (rs *Resolver) AuthLoginWithToken(w http.ResponseWriter, r *http.Request)           {}
func (rs *Resolver) AuthRecoverPassword(w http.ResponseWriter, r *http.Request)          {}
func (rs *Resolver) AuthRecoverPasswordWithToken(w http.ResponseWriter, r *http.Request) {}
func (rs *Resolver) AuthLogout(w http.ResponseWriter, r *http.Request)                   {}
