package api

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"voter-api/voter"

	"voter-api/db"

	"github.com/gin-gonic/gin"
)

type VoteAPI struct {
	voters      *db.Handler[voter.Voter]
	bootTime    time.Time
	successes   int
	badRequests int
}

func NewVoteAPI() (*VoteAPI, error) {
	dbHandler, err := db.NewHandler[voter.Voter]("voter")
	if err != nil {
		return nil, err
	}

	return &VoteAPI{voters: dbHandler, bootTime: time.Now(), successes: 0, badRequests: 0}, nil
}

func (api *VoteAPI) ListAllVoters(c *gin.Context) {
	voterList, err := api.voters.All()
	if err != nil {
		log.Println("Error Getting All Voters: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if voterList == nil {
		voterList = make([]voter.Voter, 0)
	}

	c.JSON(http.StatusOK, voterList)
	api.successes++
}

func (api *VoteAPI) DeleteAllVoters(c *gin.Context) {
	err := api.voters.Clear()
	if err != nil {
		log.Println("Error deleting All Voters: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
	api.successes++
}

func (api *VoteAPI) GetVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	voter, err := api.voters.Get(uint(id64))
	if err != nil {
		log.Println("voter not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter)
	api.successes++
}

func (api *VoteAPI) AddVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	v := voter.NewVoter()
	if err := c.ShouldBindJSON(&v); err != nil {
		log.Println("Error binding JSON: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if uint(id64) != v.VoterID {
		log.Println("URL parameter (id) does not match Request Body (voter.VoterID)")
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = api.voters.Add(v)
	if err != nil {
		log.Println("Error adding voter: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, v)
	api.successes++
}

func (api *VoteAPI) UpdateVoter(c *gin.Context) {
	var v voter.Voter
	if err := c.ShouldBindJSON(&v); err != nil {
		log.Println("Error binding JSON: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	v, err := api.voters.Update(v, voter.Voter.UpdateVoterInfo)
	if err != nil {
		log.Println("Error updating VoterInfo: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, v)
	api.successes++
}

func (api *VoteAPI) GetVoterHistory(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	v, err := api.voters.Get(uint(id64))
	if err != nil {
		log.Println("voter not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, v.GetHistory())
	api.successes++
}

func (api *VoteAPI) DeleteVoterHistory(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	v, err := api.voters.Get(uint(id64))
	if err != nil {
		log.Println("voter not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	_, err = api.voters.Update(v.DeleteHistory(), voter.Voter.UpdateVoteHistory)
	if err != nil {
		log.Println("Error updating voteHistory: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	api.successes++
}

func (api *VoteAPI) GetVoterPoll(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	voterPollIdS := c.Param("pollid")
	voterPollId64, err := strconv.ParseUint(voterPollIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	v, err := api.voters.Get(uint(voterId64))
	if err != nil {
		log.Println("voter not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	voterPoll, err := v.GetPoll(uint(voterPollId64))
	if err != nil {
		log.Println("voter poll not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voterPoll)
	api.successes++
}

func (api *VoteAPI) AddVoterPoll(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var v voter.Voter
	if err := c.ShouldBindJSON(&v); err != nil {
		log.Println("Error binding JSON: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if len(v.VoteHistory) != 1 {
		log.Println("len(v.VoteHistory) != 1")
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	vp := v.VoteHistory[0]

	v, err = api.voters.Get(uint(voterId64))
	if err != nil {
		log.Println("voter not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	updatedVoter, err := v.AddPoll(vp)
	if err != nil {
		log.Println("Error adding voter poll: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = api.voters.Update(updatedVoter, voter.Voter.UpdateVoteHistory)
	if err != nil {
		log.Println("Error updating voteHistory: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, vp)
	api.successes++
}

func (api *VoteAPI) UpdateVoterPoll(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var v voter.Voter
	if err := c.ShouldBindJSON(&v); err != nil {
		log.Println("Error binding JSON: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if len(v.VoteHistory) != 1 {
		log.Println("len(v.VoteHistory) != 1")
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	vp := v.VoteHistory[0]

	v, err = api.voters.Get(uint(voterId64))
	if err != nil {
		log.Println("voter not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	updatedVoter, err := v.UpdatePoll(vp)
	if err != nil {
		log.Println("Error updating voter poll: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = api.voters.Update(updatedVoter, voter.Voter.UpdateVoteHistory)
	if err != nil {
		log.Println("Error updating voteHistory: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, vp)
	api.successes++
}

func (api *VoteAPI) DeleteVoter(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := api.voters.Delete(uint(voterId64)); err != nil {
		log.Println("Error deleting voter: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	api.successes++
}

func (api *VoteAPI) DeleteVoterPoll(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	voterPollIdS := c.Param("pollid")
	voterPollId64, err := strconv.ParseUint(voterPollIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	v, err := api.voters.Get(uint(voterId64))
	if err != nil {
		log.Println("voter not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	updatedVoter, err := v.DeletePoll(uint(voterPollId64))
	if err != nil {
		log.Println("voter poll not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = api.voters.Update(updatedVoter, voter.Voter.UpdateVoteHistory)
	if err != nil {
		log.Println("Error updating voteHistory: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	api.successes++
}

func (api *VoteAPI) CrashSim(c *gin.Context) {
	api.badRequests++
	panic("Simulating an unexpected crash")
}

func (api *VoteAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
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
