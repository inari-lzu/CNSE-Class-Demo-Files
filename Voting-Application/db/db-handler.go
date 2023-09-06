package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
)

type Handler[T Item] struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
	keyPrefix   string
}

func NewHandler[T Item](dbName string) (*Handler[T], error) {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})

	ctx := context.Background()

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error())
		return nil, err
	}

	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	return &Handler[T]{
		cacheClient: client,
		jsonHelper:  jsonHelper,
		context:     ctx,
		keyPrefix:   dbName + ":",
	}, nil
}

//------------------------------------------------------------
// REDIS HELPERS
//------------------------------------------------------------

func (r *Handler[T]) getItemFromDB(key string) (T, error) {
	var it T
	itemObject, err := r.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return it, err
	}
	err = json.Unmarshal(itemObject.([]byte), &it)
	if err != nil {
		return it, err
	}
	return it, nil
}

func (r *Handler[T]) getKeyFromId(id uint) string {
	return fmt.Sprintf("%s%d", r.keyPrefix, id)
}

//------------------------------------------------------------
// THESE ARE THE PUBLIC FUNCTIONS THAT SUPPORT OUR TODO APP
//------------------------------------------------------------

func (r *Handler[T]) Add(it T) error {
	redisKey := r.getKeyFromId(it.GetID())
	if _, err := r.getItemFromDB(redisKey); err == nil {
		return errors.New("item already exists")
	}

	if _, err := r.jsonHelper.JSONSet(redisKey, ".", it); err != nil {
		return err
	}

	return nil
}

func (r *Handler[T]) Delete(id uint) error {
	pattern := r.getKeyFromId(id)
	numDeleted, err := r.cacheClient.Del(r.context, pattern).Result()
	if err != nil {
		return err
	}
	if numDeleted == 0 {
		return errors.New("attempted to delete non-existent item")
	}

	return nil
}

func (r *Handler[T]) Update(it T, updater func(old T, new T) (T, error)) (T, error) {
	redisKey := r.getKeyFromId(it.GetID())
	existedItem, err := r.getItemFromDB(redisKey)
	if err != nil {
		return existedItem, errors.New("item does not exist")
	}

	updatedItem, err := updater(existedItem, it)
	if err != nil {
		return existedItem, err
	}
	if _, err := r.jsonHelper.JSONSet(redisKey, ".", updatedItem); err != nil {
		return existedItem, err
	}

	return updatedItem, nil
}

func (r *Handler[T]) Get(id uint) (T, error) {
	pattern := r.getKeyFromId(id)
	existedItem, err := r.getItemFromDB(pattern)
	if err != nil {
		return existedItem, err
	}

	return existedItem, nil
}

func (r *Handler[T]) All() ([]T, error) {
	itemList := make([]T, 0)
	pattern := r.keyPrefix + "*"
	ks, _ := r.cacheClient.Keys(r.context, pattern).Result()
	for _, key := range ks {
		existedItem, err := r.getItemFromDB(key)
		if err != nil {
			return itemList, err
		}
		itemList = append(itemList, existedItem)
	}

	return itemList, nil
}

func (r *Handler[T]) Clear() error {
	pattern := r.keyPrefix + "*"
	ks, err := r.cacheClient.Keys(r.context, pattern).Result()
	if err != nil {
		return err
	}
	if len(ks) == 0 {
		return nil
	}
	numDeleted, err := r.cacheClient.Del(r.context, ks...).Result()
	if err != nil {
		return err
	}

	if numDeleted != int64(len(ks)) {
		return errors.New("one or more items could not be deleted")
	}

	return nil
}
