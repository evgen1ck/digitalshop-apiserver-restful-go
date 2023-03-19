package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/auth"
	"test-server-go/internal/queries"
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
		api_v1.RespondWithBadRequest(w, "Nickname: "+err.Error())
		return
	}
	if err := tl.Validate(email, tl.IsNotBlank(), tl.IsMinMaxLen(6, 64), tl.IsNotContainsSpace(), tl.IsEmail()); err != nil {
		api_v1.RespondWithBadRequest(w, "Email: "+err.Error())
		return
	}
	if err := tl.Validate(password, tl.IsNotBlank(), tl.IsMinMaxLen(6, 64), tl.IsNotContainsSpace()); err != nil {
		api_v1.RespondWithBadRequest(w, "Password: "+err.Error())
		return
	}

	emailDomainExists, err := tl.CheckEmailDomainExistence(email)
	if !emailDomainExists {
		api_v1.RespondWithBadRequest(w, "Email: the email domain is not exist")
		return
	}
	if err != nil {
		rs.App.Logrus.NewWarn("Error in checked the email domain: " + err.Error())
	}

	// Block 2 - checking for an existing nickname and email
	nicknameExist, emailExist, err := queries.CheckUserExistence(context.Background(), rs.App.Postgres.Pool, nickname, email)
	if err != nil {
		rs.App.Logrus.NewError("error in checked the user existence", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}
	if nicknameExist {
		api_v1.RespondWithBadRequest(w, "Nickname: this nickname is already in use")
		return
	}
	if emailExist {
		api_v1.RespondWithBadRequest(w, "Email: this email is already in use")
		return
	}

	// Block 3 - generating token and inserting a temporary account record
	randomString, err := tl.GenerateRandomString(48)
	if err != nil {
		rs.App.Logrus.NewError("error in generated random string", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}
	confirmationToken := base64.URLEncoding.EncodeToString([]byte(randomString))

	err = queries.InsertRegistrationTemp(context.Background(), rs.App.Postgres.Pool, nickname, email, password, confirmationToken)
	if err != nil {
		rs.App.Logrus.NewError("error in inserted registration temp record", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 4 - generating url and sending url on email
	_, err = tl.UrlSetParam(rs.App.Config.App.ServiceUrl+"/confirm-registration", "token", confirmationToken)
	if err != nil {
		rs.App.Logrus.NewError("error in url set param", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	//err = rs.App.Mailer.SendEmailConfirmation(nickname, email, url)
	//if err != nil {
	//	rs.App.Logrus.NewError("error in sent email confirmation", err)
	//	respondWithInternalServerError(w)
	//	return
	//}

	// Block 5 - sending the result
	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) SignupWithToken(w http.ResponseWriter, r *http.Request) {
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

	if err := tl.Validate(token, tl.IsNotBlank(), tl.IsLen(64), tl.IsNotContainsSpace()); err != nil {
		api_v1.RespondWithBadRequest(w, "Token: "+err.Error())
		return
	}

	// Block 2 - get user data and checking on exist user
	userData, err := queries.GetRegistrationTemp(context.Background(), rs.App.Postgres.Pool, token)
	if userData.Email == "" {
		api_v1.RespondWithBadRequest(w, "User not found")
		return
	}
	if err != nil {
		rs.App.Logrus.NewError("error in checked registration temp record", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 3 - hashing password and adding a user
	base64PasswordHash, base64Salt, err := auth.HashPassword(userData.Password, "")
	if err != nil {
		rs.App.Logrus.NewError("error in generated hash password", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	userUuid, err := queries.RegistrationUser(context.Background(), rs.App.Postgres.Pool, userData.Nickname, userData.Email, base64PasswordHash, base64Salt)
	if err != nil {
		rs.App.Logrus.NewError("error in registration user", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 4 - generating JWT
	jwt, err := auth.GenerateJwt(userUuid.String(), rs.App.Config.App.JwtSecret)
	if err != nil {
		rs.App.Logrus.NewError("error in generated jwt", err)
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
		Uuid:     userUuid.String(),
		Nickname: userData.Nickname,
		Email:    userData.Email,
	}
	api_v1.RespondWithCreated(w, response)
}
