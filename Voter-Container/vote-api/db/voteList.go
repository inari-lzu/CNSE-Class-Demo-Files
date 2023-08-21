package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

type voterPoll struct {
	PollID   uint      `json:"id"`
	VoteDate time.Time `json:"date"`
}

type Voter struct {
	VoterID     uint        `json:"id"`
	FirstName   string      `json:"firstName"`
	LastName    string      `json:"lastName"`
	VoteHistory []voterPoll `json:"history"`
}

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "voter:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

type VoterList struct {
	cache
}

func NewVoterList() (*VoterList, error) {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	return NewWithCacheInstance(redisUrl)
}

func NewWithCacheInstance(location string) (*VoterList, error) {

	//Connect to redis.  Other options can be provided, but the
	//defaults are OK
	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	//We use this context to coordinate betwen our go code and
	//the redis operaitons
	ctx := context.Background()

	//This is the reccomended way to ensure that our redis connection
	//is working
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error())
		return nil, err
	}

	//By default, redis manages keys and values, where the values
	//are either strings, sets, maps, etc.  Redis has an extension
	//module called ReJSON that allows us to store JSON objects
	//however, we need a companion library in order to work with it
	//Below we create an instance of the JSON helper and associate
	//it with our redis connnection
	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	//Return a pointer to a new ToDo struct
	return &VoterList{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}, nil
}

//------------------------------------------------------------
// REDIS HELPERS
//------------------------------------------------------------

// func isRedisNilError(err error) bool {
// 	return errors.Is(err, redis.Nil) || err.Error() == RedisNilError
// }

func (v *VoterList) getItemFromRedis(key string, voter *Voter) error {
	itemObject, err := v.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	err = json.Unmarshal(itemObject.([]byte), voter)
	if err != nil {
		return err
	}

	return nil
}

func redisKeyFromId(id uint) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

//------------------------------------------------------------
// THESE ARE THE PUBLIC FUNCTIONS THAT SUPPORT OUR TODO APP
//------------------------------------------------------------

func (v *VoterList) AddVoter(voter Voter) error {
	redisKey := redisKeyFromId(voter.VoterID)
	var existingVoter Voter
	if err := v.getItemFromRedis(redisKey, &existingVoter); err == nil {
		return errors.New("voter already exists")
	}

	if _, err := v.jsonHelper.JSONSet(redisKey, ".", voter); err != nil {
		return err
	}

	return nil
}

func (v *VoterList) DeleteVoter(id uint) error {
	pattern := redisKeyFromId(id)
	numDeleted, err := v.cacheClient.Del(v.context, pattern).Result()
	if err != nil {
		return err
	}
	if numDeleted == 0 {
		return errors.New("attempted to delete non-existent voter")
	}

	return nil
}

func (v *VoterList) DeleteAllVoters() error {
	pattern := RedisKeyPrefix + "*"
	ks, err := v.cacheClient.Keys(v.context, pattern).Result()
	if err != nil {
		return err
	}
	if len(ks) == 0 {
		return nil
	}
	numDeleted, err := v.cacheClient.Del(v.context, ks...).Result()
	if err != nil {
		return err
	}

	if numDeleted != int64(len(ks)) {
		return errors.New("one or more items could not be deleted")
	}

	return nil
}

func (v *VoterList) UpdateVoter(voter Voter) (Voter, error) {
	redisKey := redisKeyFromId(voter.VoterID)
	var existingVoter Voter
	if err := v.getItemFromRedis(redisKey, &existingVoter); err != nil {
		return Voter{}, errors.New("item does not exist")
	}
	voter.VoteHistory = existingVoter.VoteHistory // we don't expect to change voteHistory
	if _, err := v.jsonHelper.JSONSet(redisKey, ".", voter); err != nil {
		return Voter{}, err
	}

	return voter, nil
}

func (v *VoterList) GetVoter(id uint) (Voter, error) {
	var voter Voter
	pattern := redisKeyFromId(id)
	err := v.getItemFromRedis(pattern, &voter)
	if err != nil {
		return Voter{}, err
	}

	return voter, nil
}

func (v *VoterList) GetAllVoters() ([]Voter, error) {
	var voterList []Voter
	var voter Voter

	pattern := RedisKeyPrefix + "*"
	ks, _ := v.cacheClient.Keys(v.context, pattern).Result()
	for _, key := range ks {
		err := v.getItemFromRedis(key, &voter)
		if err != nil {
			return nil, err
		}
		voterList = append(voterList, voter)
	}

	return voterList, nil
}

func (v *VoterList) AddVoterPoll(voterId uint, newvp voterPoll) error {
	voter, err := v.GetVoter(voterId)
	if err != nil {
		return errors.New("voter does not exist")
	}
	for _, vp := range voter.VoteHistory {
		if newvp.PollID == vp.PollID {
			return errors.New("voter poll already exists")
		}
	}
	voter.VoteHistory = append(voter.VoteHistory, newvp)
	_, err = v.jsonHelper.JSONSet(redisKeyFromId(voter.VoterID), ".", voter)
	if err != nil {
		return errors.New("update redis data fail")
	}
	return nil
}

func (v *VoterList) DeleteVoterPoll(voterId uint, pollid uint) error {
	voter, err := v.GetVoter(voterId)
	if err != nil {
		return errors.New("voter does not exist")
	}
	for i, vp := range voter.VoteHistory {
		if vp.PollID == pollid {
			voter.VoteHistory = append(voter.VoteHistory[:i], voter.VoteHistory[i+1:]...)
			_, err := v.jsonHelper.JSONSet(redisKeyFromId(voter.VoterID), ".", voter)
			if err != nil {
				return errors.New("update redis data fail")
			}
			return nil
		}
	}
	return errors.New("voter poll does not exist")
}

func (v *VoterList) DeleteVoteHistory(voterId uint) error {
	voter, err := v.GetVoter(voterId)
	if err != nil {
		return errors.New("voter does not exist")
	}
	voter.VoteHistory = []voterPoll{}
	_, err = v.jsonHelper.JSONSet(redisKeyFromId(voter.VoterID), ".", voter)
	if err != nil {
		return errors.New("update redis data fail")
	}
	return nil
}

func (v *VoterList) UpdateVoterPoll(voterId uint, newvp voterPoll) error {
	voter, err := v.GetVoter(voterId)
	if err != nil {
		return errors.New("voter does not exist")
	}
	pollid := newvp.PollID
	for i, vp := range voter.VoteHistory {
		if vp.PollID == pollid {
			voter.VoteHistory[i] = newvp

			_, err := v.jsonHelper.JSONSet(redisKeyFromId(voter.VoterID), ".", voter)
			if err != nil {
				return errors.New("update redis data fail")
			}
			return nil
		}
	}
	return errors.New("voter poll does not exist")
}

func (v *VoterList) GetVoterPoll(voterId uint, pollid uint) (voterPoll, error) {
	voter, err := v.GetVoter(voterId)
	if err != nil {
		return voterPoll{}, errors.New("voter does not exist")
	}
	for _, vp := range voter.VoteHistory {
		if vp.PollID == pollid {
			return vp, nil
		}
	}
	return voterPoll{}, errors.New("voter poll does not exist")
}

func (v *VoterList) GetVoteHistory(voterId uint) ([]voterPoll, error) {
	voter, err := v.GetVoter(voterId)
	if err != nil {
		return []voterPoll{}, errors.New("voter does not exist")
	}
	return voter.VoteHistory, nil
}

//------------------------------------------------------------
// Helpers for debugging
//------------------------------------------------------------

func (v *VoterList) PrintVoter(voter Voter) {
	jsonBytes, _ := json.MarshalIndent(voter, "", "  ")
	fmt.Println(string(jsonBytes))
}

func (v *VoterList) PrintAllVoters(vl []Voter) {
	for _, voter := range vl {
		v.PrintVoter(voter)
	}
}

func (v *VoterList) JsonToVoter(jsonString string) (Voter, error) {
	var voter Voter
	err := json.Unmarshal([]byte(jsonString), &voter)
	if err != nil {
		return Voter{}, err
	}

	return voter, nil
}
