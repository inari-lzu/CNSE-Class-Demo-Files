package api

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"votes-api/vote"

	"db"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type VoteAPI struct {
	votes            *db.Handler[vote.Vote]
	bootTime         time.Time
	successes        int
	badRequests      int
	apiClient        *resty.Client
	hostName         string
	voterApiInternal string
	pollApiInternal  string
	voterApiExternal string
	pollApiExternal  string
}

func NewVotesAPI() (*VoteAPI, error) {
	hostName := os.Getenv("HOST_NAME")
	voterApiInternal := os.Getenv("VOTER_API_INTERNAL")
	pollApiInternal := os.Getenv("POLL_API_INTERNAL")
	voterApiExternal := os.Getenv("VOTER_API_EXTERAL")
	pollApiExternal := os.Getenv("POLL_API_EXTERAL")

	log.Println("HOST_NAME: " + hostName)
	log.Println("VOTER_API_INTERNAL: " + voterApiInternal)
	log.Println("POLL_API_INTERNAL: " + pollApiInternal)
	log.Println("VOTER_API_EXTERAL: " + voterApiExternal)
	log.Println("POLL_API_EXTERAL: " + pollApiExternal)

	dbHandler, err := db.NewHandler[vote.Vote]("vote")
	if err != nil {
		return nil, err
	}

	return &VoteAPI{
		votes:            dbHandler,
		bootTime:         time.Now(),
		successes:        0,
		badRequests:      0,
		apiClient:        resty.New(),
		hostName:         hostName,
		voterApiInternal: voterApiInternal,
		pollApiInternal:  pollApiInternal,
		voterApiExternal: voterApiExternal,
		pollApiExternal:  pollApiExternal,
	}, nil
}

func (api *VoteAPI) ListAllVotes(c *gin.Context) {
	voteList, err := api.votes.All()
	if err != nil {
		log.Println("Error Getting All Votes: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if voteList == nil {
		voteList = make([]vote.Vote, 0)
	}

	urlList := make([]string, 0)
	for _, v := range voteList {
		urlList = append(urlList, v.ToLinks(api.hostName, api.voterApiExternal, api.pollApiExternal).Voter)
	}

	c.JSON(http.StatusOK, urlList)
	api.successes++
}

func (api *VoteAPI) DeleteAllVotes(c *gin.Context) {
	// delete all associated vote history
	voteList, err := api.votes.All()
	if err != nil {
		log.Println("Error Getting All Votes: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, v := range voteList {
		err = api.safeDeleteVote(v, c)
		if err != nil {
			api.badRequests++
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Some of the votes are not deleted success"})
			return
		}
	}

	if err != nil {
		log.Println("Error deleting All Votes: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
	api.successes++
}

func (api *VoteAPI) GetVote(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	v, err := api.votes.Get(uint(id64))
	if err != nil {
		log.Println("vote not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, v.ToLinks(api.hostName, api.voterApiExternal, api.pollApiExternal))
	api.successes++
}

func (api *VoteAPI) validateVoteLinks(links vote.Links) error {
	checkVoterExist, err := api.apiClient.R().Get(links.Voter)
	if err != nil || checkVoterExist.StatusCode() != 200 {
		return errors.New("associated Voter doesn't existed")
	}

	checkPollExist, err := api.apiClient.R().Get(links.Poll)
	if err != nil || checkPollExist.StatusCode() != 200 {
		return errors.New("associated Poll doesn't existed")
	}

	checkOptionExist, err := api.apiClient.R().Get(links.Choice)
	if err != nil || checkOptionExist.StatusCode() != 200 {
		return errors.New("associated vote option doesn't existed")
	}

	return nil
}

func (api *VoteAPI) AddVote(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	v := vote.NewVote()
	if err := c.ShouldBindJSON(&v); err != nil {
		log.Println("Error binding JSON: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if uint(id64) != v.VoteID {
		log.Println("URL parameter (id) does not match Request Body (v.VoteID)")
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// exam whether associate entities existed
	links := v.ToLinks(api.hostName, api.voterApiInternal, api.pollApiInternal)
	err = api.validateVoteLinks(links)
	if err != nil {
		api.badRequests++
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add new vote into redis
	err = api.votes.Add(v)
	if err != nil {
		log.Println("Error adding vote: ", err)
		api.badRequests++
		emsg := "Error adding vote: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": emsg})
		return
	}

	// try call voter-api to add voterPoll to associate voter's voting history
	addVoterPoll, err := api.apiClient.R().
		SetBody(v.ToVoteHistoryRecord()).
		Post(links.VoterPoll)
	if err != nil || addVoterPoll.StatusCode() != 200 {
		api.votes.Delete(v.VoteID) // undo add vote
		api.badRequests++
		emsg := "Adding voterPoll to associate voter's voting history fail. One voter can only has one voting in a poll."
		c.JSON(http.StatusInternalServerError, gin.H{"error": emsg})
		return
	}

	c.JSON(http.StatusOK, v.ToLinks(api.hostName, api.voterApiExternal, api.pollApiExternal))
	api.successes++
}

func (api *VoteAPI) UpdateVote(c *gin.Context) {
	var v vote.Vote
	if err := c.ShouldBindJSON(&v); err != nil {
		log.Println("Error binding JSON: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// exam whether associate entities existed
	links := v.ToLinks(api.hostName, api.voterApiInternal, api.pollApiInternal)
	err := api.validateVoteLinks(links)
	if err != nil {
		api.badRequests++
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	prev, err := api.votes.Get(v.VoteID)
	if err != nil {
		log.Println("Vote not exist: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	v, err = api.votes.Update(v, vote.Vote.Update)
	if err != nil {
		log.Println("Error updating vote: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// try call voter-api to update voterPoll to associate voter's voting history
	updateVoterPoll, err := api.apiClient.R().
		SetBody(v.ToVoteHistoryRecord()).
		Put(links.VoterPoll)
	if err != nil || updateVoterPoll.StatusCode() != 200 {
		api.votes.Update(prev, vote.Vote.Update) // undo update
		api.badRequests++
		emsg := "Updating voterPoll to associate voter's voting history fail"
		c.JSON(http.StatusInternalServerError, gin.H{"error": emsg})
		return
	}

	c.JSON(http.StatusOK, v.ToLinks(api.hostName, api.voterApiExternal, api.pollApiExternal))
	api.successes++
}

func (api *VoteAPI) DeleteVote(c *gin.Context) {
	voteIdS := c.Param("id")
	voteId64, err := strconv.ParseUint(voteIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	voteId := uint(voteId64)

	v, err := api.votes.Get(voteId)
	if err != nil {
		log.Println("Vote not exist: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = api.safeDeleteVote(v, c)
	if err != nil {
		api.badRequests++
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
	api.successes++
}

func (api *VoteAPI) safeDeleteVote(v vote.Vote, c *gin.Context) error {
	if err := api.votes.Delete(v.VoteID); err != nil {
		return errors.New("Error deleting vote: " + err.Error())

	}

	// try call voter-api to delete voterPoll from associate voter's voting history
	links := v.ToLinks(api.hostName, api.voterApiInternal, api.pollApiInternal)
	deleteVoterPoll, err := api.apiClient.R().
		SetBody(v.ToVoteHistoryRecord()).
		Delete(links.VoterPoll)
	if err != nil || deleteVoterPoll.StatusCode() != 200 {
		api.votes.Add(v) // undo delete
		return errors.New("deleting voterPoll from associate voter's voting history fail")
	}
	return nil
}

func (api *VoteAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"api_name":                    "votes-api",
			"status":                      http.StatusOK,
			"gin_version":                 gin.Version,
			"api_uptime":                  time.Since(api.bootTime).String(),
			"total_api_calls":             api.successes + api.badRequests,
			"total_api_calls_succeed":     api.successes,
			"total_api_calls_with_errors": api.badRequests,
		},
	)
	api.successes++
}
