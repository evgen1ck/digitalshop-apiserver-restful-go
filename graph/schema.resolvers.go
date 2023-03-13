package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"test-server-go/graph/model"
	"test-server-go/internal/auth"
	"test-server-go/internal/queries"
	tl "test-server-go/internal/tools"
)

// AuthSignup is the resolver for the authSignup field.
func (r *mutationResolver) AuthSignup(ctx context.Context, input model.SignupInput) (bool, error) {
	// Block 1 - data validation
	nickname := strings.TrimSpace(input.Nickname)
	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := strings.TrimSpace(input.Password)

	if err := tl.Validate(nickname, tl.IsMinMaxLen(5, 32), tl.IsNotContainsSpace(), tl.IsNickname()); err != nil {
		return false, errors.New("nickname: " + err.Error())
	}
	if err := tl.Validate(email, tl.IsMinMaxLen(6, 64), tl.IsNotContainsSpace(), tl.IsEmail()); err != nil {
		return false, errors.New("email: " + err.Error())
	}
	if err := tl.Validate(password, tl.IsMinMaxLen(6, 64), tl.IsNotContainsSpace()); err != nil {
		return false, errors.New("password: " + err.Error())
	}
	emailDomainExists, err := tl.CheckEmailDomainExistence(email)
	if !emailDomainExists {
		return false, errors.New("email: the email domain is not exist")
	}
	if err != nil {
		r.App.Logrus.NewWarn("error in checked the email domain: " + err.Error())
	}

	// Block 2 - checking for an existing nickname and email
	nicknameExist, emailExist, err := queries.CheckUserExistence(ctx, r.App.Postgres.Pool, nickname, email)
	if err != nil {
		r.App.Logrus.NewError("error in checked the user existence", err)
		return false, errors.New("system error")
	}
	if nicknameExist {
		return false, errors.New("nickname: this nickname is already in use")
	}
	if emailExist {
		return false, errors.New("email: this email is already in use")
	}

	// Block 3 - generating token and inserting a temporary account record
	randomString, err := tl.GenerateRandomString(48)
	if err != nil {
		r.App.Logrus.NewError("error in generated random string", err)
		return false, errors.New("system error")
	}
	confirmationToken := base64.URLEncoding.EncodeToString([]byte(randomString))

	err = queries.InsertRegistrationTemp(ctx, r.App.Postgres.Pool, nickname, email, password, confirmationToken)
	if err != nil {
		r.App.Logrus.NewError("error in inserted registration temp record", err)
		return false, errors.New("system error")
	}

	// Block 4 - generating url and sending url on email
	url, err := tl.UrlSetParam("https://digitalshop.evgenick.com/confirm-registration", "token", confirmationToken)
	if err != nil {
		r.App.Logrus.NewError("error in url set param", err)
		return false, errors.New("system error")
	}

	err = r.App.Mailer.SendEmailConfirmation(nickname, email, url)
	if err != nil {
		r.App.Logrus.NewError("error in sent email confirmation", err)
		return false, errors.New("system error")
	}

	// Block 5 - sending the result
	return true, nil
}

// AuthSignupWithToken is the resolver for the authSignupWithToken field.
func (r *mutationResolver) AuthSignupWithToken(ctx context.Context, token string) (*model.AuthUserPayload, error) {
	// Block 1 - data validation
	token = strings.TrimSpace(token)

	if err := tl.Validate(token, tl.IsLen(64), tl.IsNotContainsSpace()); err != nil {
		return nil, errors.New("token: " + err.Error())
	}

	// Block 2 - get user data and checking on exist user
	userData, err := queries.GetRegistrationTemp(ctx, r.App.Postgres.Pool, token)
	if userData.Email == "" {
		return nil, errors.New("user not found")
	}
	if err != nil {
		r.App.Logrus.NewError("error in checked registration temp record", err)
		return nil, errors.New("system error")
	}

	// Block 3 - hashing password and adding a user
	base64PasswordHash, base64Salt, err := auth.HashPassword(userData.Password, "")
	if err != nil {
		r.App.Logrus.NewError("error in generated hash password", err)
		return nil, errors.New("system error")
	}

	userUuid, err := queries.RegistrationUser(ctx, r.App.Postgres.Pool, userData.Nickname, userData.Email, base64PasswordHash, base64Salt)
	if err != nil {
		r.App.Logrus.NewError("error in registration user", err)
		return nil, errors.New("system error")
	}

	// Block 4 - generating JWT
	jwt, err := auth.GenerateJwt(userUuid.String(), r.App.Config.App.JwtSecret)
	if err != nil {
		r.App.Logrus.NewError("error in generated jwt", err)
		return nil, errors.New("system error")
	}

	// Block 5 - sending the result
	return &model.AuthUserPayload{
		Token:    jwt,
		UUID:     userUuid.String(),
		Nickname: userData.Nickname,
		Email:    userData.Email,
	}, nil
}

// AuthLogin is the resolver for the authLogin field.
func (r *mutationResolver) AuthLogin(ctx context.Context, input model.LoginInput) (*model.AuthUserPayload, error) {
	panic(fmt.Errorf("not implemented: AuthLogin - authLogin"))
}

// AuthLogout is the resolver for the authLogout field.
func (r *mutationResolver) AuthLogout(ctx context.Context, token string) (bool, error) {
	panic(fmt.Errorf("not implemented: AuthLogout - authLogout"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
