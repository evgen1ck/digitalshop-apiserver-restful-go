package api_v1

//import (
//	"context"
//	"encoding/base64"
//	"encoding/json"
//	"net/http"
//	"os"
//	"strings"
//	"test-server-go/internal/auth"
//	"test-server-go/internal/queries"
//	tl "test-server-go/internal/tools"
//)
//
//type Album struct {
//	ID     int64   `json:"id"`
//	Title  string  `json:"title"`
//	Artist string  `json:"artist"`
//	Price  float64 `json:"price"`
//}
//
//func (rh *RouteHandler) AuthSignup(w http.ResponseWriter, r *http.Request) {
//	var input struct {
//		Nickname string `json:"nickname"`
//		Email    string `json:"email"`
//		Password string `json:"password"`
//	}
//	decodeErr := json.NewDecoder(r.Body).Decode(&input)
//	if decodeErr != nil {
//		RespondWithBadRequest(w, "")
//		return
//	}
//
//	// Block 1 - data validation
//	nickname := strings.TrimSpace(input.Nickname)
//	email := strings.TrimSpace(strings.ToLower(input.Email))
//	password := strings.TrimSpace(input.Password)
//
//	if err := tl.Validate(nickname, tl.IsNotBlank(), tl.IsMinMaxLen(5, 32), tl.IsNotContainsSpace(), tl.IsNickname()); err != nil {
//		RespondWithBadRequest(w, "Nickname: "+err.Error())
//		return
//	}
//	if err := tl.Validate(email, tl.IsNotBlank(), tl.IsMinMaxLen(6, 64), tl.IsNotContainsSpace(), tl.IsEmail()); err != nil {
//		RespondWithBadRequest(w, "Email: "+err.Error())
//		return
//	}
//	if err := tl.Validate(password, tl.IsNotBlank(), tl.IsMinMaxLen(6, 64), tl.IsNotContainsSpace()); err != nil {
//		RespondWithBadRequest(w, "Password: "+err.Error())
//		return
//	}
//
//	emailDomainExists, err := tl.CheckEmailDomainExistence(email)
//	if !emailDomainExists {
//		RespondWithBadRequest(w, "Email: the email domain is not exist")
//		return
//	}
//	if err != nil {
//		rh.App.Logrus.NewWarn("Error in checked the email domain: " + err.Error())
//	}
//
//	// Block 2 - checking for an existing nickname and email
//	nicknameExist, emailExist, err := queries.CheckUserExistence(context.Background(), rh.App.Postgres.Pool, nickname, email)
//	if err != nil {
//		rh.App.Logrus.NewError("error in checked the user existence", err)
//		RespondWithInternalServerError(w)
//		return
//	}
//	if nicknameExist {
//		RespondWithBadRequest(w, "Nickname: this nickname is already in use")
//		return
//	}
//	if emailExist {
//		RespondWithBadRequest(w, "Email: this email is already in use")
//		return
//	}
//
//	// Block 3 - generating token and inserting a temporary account record
//	randomString, err := tl.GenerateRandomString(48)
//	if err != nil {
//		rh.App.Logrus.NewError("error in generated random string", err)
//		RespondWithInternalServerError(w)
//		return
//	}
//	confirmationToken := base64.URLEncoding.EncodeToString([]byte(randomString))
//
//	err = queries.InsertRegistrationTemp(context.Background(), rh.App.Postgres.Pool, nickname, email, password, confirmationToken)
//	if err != nil {
//		rh.App.Logrus.NewError("error in inserted registration temp record", err)
//		RespondWithInternalServerError(w)
//		return
//	}
//
//	// Block 4 - generating url and sending url on email
//	_, err = tl.UrlSetParam(rh.App.Config.App.ServiceUrl+"/confirm-registration", "token", confirmationToken)
//	if err != nil {
//		rh.App.Logrus.NewError("error in url set param", err)
//		RespondWithInternalServerError(w)
//		return
//	}
//
//	//err = rh.App.Mailer.SendEmailConfirmation(nickname, email, url)
//	//if err != nil {
//	//	rh.App.Logrus.NewError("error in sent email confirmation", err)
//	//	RespondWithInternalServerError(w)
//	//	return
//	//}
//
//	// Block 5 - sending the result
//	w.WriteHeader(http.StatusNoContent)
//}
//
//func (rh *RouteHandler) SignupWithToken(w http.ResponseWriter, r *http.Request) {
//	var input struct {
//		Token string `json:"token"`
//	}
//	decodeErr := json.NewDecoder(r.Body).Decode(&input)
//	if decodeErr != nil {
//		RespondWithBadRequest(w, "")
//		return
//	}
//
//	// Block 1 - data validation
//	token := strings.TrimSpace(input.Token)
//
//	if err := tl.Validate(token, tl.IsNotBlank(), tl.IsLen(64), tl.IsNotContainsSpace()); err != nil {
//		RespondWithBadRequest(w, "Token: "+err.Error())
//		return
//	}
//
//	// Block 2 - get user data and checking on exist user
//	userData, err := queries.GetRegistrationTemp(context.Background(), rh.App.Postgres.Pool, token)
//	if userData.Email == "" {
//		RespondWithBadRequest(w, "User not found")
//		return
//	}
//	if err != nil {
//		rh.App.Logrus.NewError("error in checked registration temp record", err)
//		RespondWithInternalServerError(w)
//		return
//	}
//
//	// Block 3 - hashing password and adding a user
//	base64PasswordHash, base64Salt, err := auth.HashPassword(userData.Password, "")
//	if err != nil {
//		rh.App.Logrus.NewError("error in generated hash password", err)
//		RespondWithInternalServerError(w)
//		return
//	}
//
//	userUuid, err := queries.RegistrationUser(context.Background(), rh.App.Postgres.Pool, userData.Nickname, userData.Email, base64PasswordHash, base64Salt)
//	if err != nil {
//		rh.App.Logrus.NewError("error in registration user", err)
//		RespondWithInternalServerError(w)
//		return
//	}
//
//	// Block 4 - generating JWT
//	jwt, err := auth.GenerateJwt(userUuid.String(), rh.App.Config.App.JwtSecret)
//	if err != nil {
//		rh.App.Logrus.NewError("error in generated jwt", err)
//		RespondWithInternalServerError(w)
//		return
//	}
//
//	// Block 5 - sending the result
//	response := struct {
//		Token    string `json:"token"`
//		Uuid     string `json:"uuid"`
//		Nickname string `json:"nickname"`
//		Email    string `json:"email"`
//	}{
//		Token:    jwt,
//		Uuid:     userUuid.String(),
//		Nickname: userData.Nickname,
//		Email:    userData.Email,
//	}
//	RespondWithCreated(w, response)
//}
//
//func getAllAlbums(w http.ResponseWriter, r *http.Request) {
//	var albums = []Album{
//		{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
//		{ID: 2, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
//		{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
//	}
//
//	ipAddress := r.RemoteAddr
//
//	err := json.NewEncoder(w).Encode(map[string]interface{}{
//		"ip_address": ipAddress,
//		"users":      albums,
//	})
//	if err != nil {
//		os.Exit(1)
//	}
//}
