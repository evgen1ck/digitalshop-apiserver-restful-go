package storage

import (
	"context"
	"github.com/redis/go-redis/v9"
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
	exists, err := rdb.Client.Exists(ctx, BlockedTokenPath+token).Result()
	if err != nil {
		return err
	} else if exists > 0 {
		return QueryExists
	}

	if err := rdb.Client.Set(ctx, BlockedTokenPath+token, "true", expiration).Err(); err != nil {
		return err
	}
	return nil
}

func CheckBlockedTokenExists(ctx context.Context, rdb *Redis, token string) (bool, error) {
	result, err := rdb.Client.Get(ctx, BlockedTokenPath+token).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return result == "true", nil
	}
}

func CreateTempRegistration(ctx context.Context, rdb *Redis, nickname, email, password, confirmationToken string, expiration time.Duration) error {
	exists, err := rdb.Client.Exists(ctx, TempRegistrationPath+confirmationToken).Result()
	if err != nil {
		return err
	} else if exists > 0 {
		return QueryExists
	}

	if err = execInPipeline(ctx, rdb.Client, func(pipe redis.Pipeliner) error {
		if err := pipe.HSet(ctx, TempRegistrationPath+confirmationToken, "nickname", nickname, "email", email, "password", password).Err(); err != nil {
			return err
		}

		if err = pipe.Expire(ctx, TempRegistrationPath+confirmationToken, expiration).Err(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
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

	if err = rdb.Client.Del(ctx, TempRegistrationPath+confirmationToken).Err(); err != nil {
		return err
	}

	return nil
}
