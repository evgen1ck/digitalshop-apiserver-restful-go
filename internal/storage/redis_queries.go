package storage

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

// Names style:
// For creating a record: Create<Type>
// For checking the existence of a record: Check<Type>Exists
// For getting a record: Get<Type>
// For updating a record: Update<Type>
// For deleting a record: Delete<Type>

const (
	BlockedTokenPath     = "jwt_stoplist:"
	TempRegistrationPath = "registration_temp_data:"
)

func CreateBlockedToken(ctx context.Context, rdb *Redis, token string, expiration time.Duration) error {
	token = strings.ToLower(token)

	exists, err := rdb.Client.Exists(ctx, BlockedTokenPath+token).Result()
	if err != nil {
		return err
	} else if exists > 0 {
		return QueryExists
	}

	err = rdb.Client.Set(ctx, BlockedTokenPath+token, "true", expiration).Err()
	return err
}

func CheckBlockedTokenExists(ctx context.Context, rdb *Redis, token string) (bool, error) {
	token = strings.ToLower(token)

	result, err := rdb.Client.Get(ctx, BlockedTokenPath+token).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return result == "true", err
	}
}

func CreateTempRegistration(ctx context.Context, rdb *Redis, nickname, email, password, confirmationToken string, expiration time.Duration) error {
	nickname = strings.ToLower(nickname)
	email = strings.ToLower(email)

	exists, err := rdb.Client.Exists(ctx, TempRegistrationPath+confirmationToken).Result()
	if err != nil {
		return err
	} else if exists > 0 {
		return QueryExists
	}

	err = execInPipeline(ctx, rdb.Client, func(pipe redis.Pipeliner) error {
		if err = pipe.HSet(ctx, TempRegistrationPath+confirmationToken, "nickname", nickname, "email", email, "password", password).Err(); err != nil {
			return err
		}

		err = pipe.Expire(ctx, TempRegistrationPath+confirmationToken, expiration).Err()
		return err
	})

	return err
}

func GetTempRegistration(ctx context.Context, rdb *Redis, confirmationToken string) (string, string, string, error) {
	var nickname, email, password string

	data, err := rdb.Client.HGetAll(ctx, TempRegistrationPath+confirmationToken).Result()
	if err != nil {
		return nickname, email, password, err
	}

	if len(data) == 0 {
		return nickname, email, password, NoResults
	}

	nickname = data["nickname"]
	email = data["email"]
	password = data["password"]

	return nickname, email, password, nil
}

func DeleteTempRegistration(ctx context.Context, rdb *Redis, confirmationToken string) error {
	exists, err := rdb.Client.Exists(ctx, TempRegistrationPath+confirmationToken).Result()
	if err != nil {
		return err
	} else if exists == 0 {
		return NoResults
	}

	err = rdb.Client.Del(ctx, TempRegistrationPath+confirmationToken).Err()

	return err
}
