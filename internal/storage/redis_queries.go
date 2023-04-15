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

func CreateBlockedToken(ctx context.Context, rdb *Redis, token string, ttl time.Duration) error {
	err := rdb.Client.Set(ctx, "blacklist:"+token, "true", ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func CheckBlockedTokenExists(ctx context.Context, rdb *Redis, token string) (bool, error) {
	result, err := rdb.Client.Get(ctx, "blacklist:"+token).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return result == "true", nil
	}
}
