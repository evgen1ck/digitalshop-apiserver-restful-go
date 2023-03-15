package api

import (
	"encoding/json"
	"net/http"
	"os"
)

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type SignupInput struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (rh *RouteHandler) authSignup(w http.ResponseWriter, r *http.Request) {
	//var input SignupInput
	//decodeErr := json.NewDecoder(r.Body).Decode(&input)
	//if decodeErr != nil {
	//	respondWithBadRequest(w, "Invalid request payload")
	//}
	//
	//// Block 1 - data validation
	//nickname := strings.TrimSpace(input.Nickname)
	//email := strings.TrimSpace(strings.ToLower(input.Email))
	//password := strings.TrimSpace(input.Password)
	//
	//if err := tl.Validate(nickname, tl.IsMinMaxLen(5, 32), tl.IsNotContainsSpace(), tl.IsNickname()); err != nil {
	//	respondWithBadRequest(w, "Nickname: "+err.Error())
	//}
	//if err := tl.Validate(email, tl.IsMinMaxLen(6, 64), tl.IsNotContainsSpace(), tl.IsEmail()); err != nil {
	//	respondWithBadRequest(w, "Email: "+err.Error())
	//}
	//if err := tl.Validate(password, tl.IsMinMaxLen(6, 64), tl.IsNotContainsSpace()); err != nil {
	//	respondWithBadRequest(w, "Password: "+err.Error())
	//}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"ip_address": rh.App.Config.App.ServiceName,
	})

	//emailDomainExists, err := tl.CheckEmailDomainExistence(email)
	//if !emailDomainExists {
	//	respondWithBadRequest(w, "Email: the email domain is not exist")
	//}
	//if err != nil {
	//	r.Server.Logrus.NewWarn("error in checked the email domain: " + err.Error())
	//}
	//
	//// Block 2 - checking for an existing nickname and email
	//nicknameExist, emailExist, err := queries.CheckUserExistence(ctx, r.App.Postgres.Pool, nickname, email)
	//if err != nil {
	//	r.App.Logrus.NewError("error in checked the user existence", err)
	//	return false, errors.New("system error")
	//}
	//if nicknameExist {
	//	return false, errors.New("nickname: this nickname is already in use")
	//}
	//if emailExist {
	//	return false, errors.New("email: this email is already in use")
	//}
	//
	//// Block 3 - generating token and inserting a temporary account record
	//randomString, err := tl.GenerateRandomString(48)
	//if err != nil {
	//	r.App.Logrus.NewError("error in generated random string", err)
	//	return false, errors.New("system error")
	//}
	//confirmationToken := base64.URLEncoding.EncodeToString([]byte(randomString))
	//
	//err = queries.InsertRegistrationTemp(ctx, r.App.Postgres.Pool, nickname, email, password, confirmationToken)
	//if err != nil {
	//	r.App.Logrus.NewError("error in inserted registration temp record", err)
	//	return false, errors.New("system error")
	//}
	//
	//// Block 4 - generating url and sending url on email
	//url, err := tl.UrlSetParam("https://digitalshop.evgenick.com/confirm-registration", "token", confirmationToken)
	//if err != nil {
	//	r.App.Logrus.NewError("error in url set param", err)
	//	return false, errors.New("system error")
	//}
	//
	//err = r.App.Mailer.SendEmailConfirmation(nickname, email, url)
	//if err != nil {
	//	r.App.Logrus.NewError("error in sent email confirmation", err)
	//	return false, errors.New("system error")
	//}
	//
	//// Block 5 - sending the result
	//return true, nil

}

func getAllAlbums(w http.ResponseWriter, r *http.Request) {
	var albums = []Album{
		{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: 2, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}

	ipAddress := r.RemoteAddr

	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"ip_address": ipAddress,
		"users":      albums,
	})
	if err != nil {
		os.Exit(1)
	}
}
