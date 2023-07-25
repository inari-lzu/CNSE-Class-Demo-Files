package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"drexel.edu/todo-api/db"
	"github.com/gin-gonic/gin"
)

type VoteAPI struct {
	db *db.VoterList
}

func NewVoteAPI() (*VoteAPI, error) {
	dbHandler, err := db.NewVoterList()
	if err != nil {
		return nil, err
	}

	return &VoteAPI{db: dbHandler}, nil
}

func (va *VoteAPI) ListAllVoters(c *gin.Context) {

	voterList, err := va.db.GetAllVoters()
	if err != nil {
		log.Println("Error Getting All Voters: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if voterList == nil {
		voterList = make([]db.Voter, 0)
	}

	c.JSON(http.StatusOK, voterList)
}

func (va *VoteAPI) GetVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	voter, err := va.db.GetVoter(uint(id64))
	if err != nil {
		log.Println("voter not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter)
}

func (va *VoteAPI) AddVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var voter db.Voter
	if err := c.ShouldBindJSON(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if uint(id64) != voter.VoterID {
		log.Println("URL parameter (id) does not match Request Body (voter.VoterID)")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := va.db.AddVoter(voter); err != nil {
		log.Println("Error adding voter: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, voter)
}

func (va *VoteAPI) GetVoterHistory(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	voter, err := va.db.GetVoter(uint(id64))
	if err != nil {
		log.Println("voter not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter.VoteHistory)
}

func (va *VoteAPI) GetVoterPoll(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	voterPollIdS := c.Param("pollid")
	voterPollId64, err := strconv.ParseUint(voterPollIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voterPoll, err := va.db.GetVoterPoll(uint(voterId64), uint(voterPollId64))
	if err != nil {
		log.Println("voter poll not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voterPoll)
}

func (va *VoteAPI) AddVoterPoll(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voterPollIdS := c.Param("pollid")
	voterPollId64, err := strconv.ParseUint(voterPollIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pollId := uint(voterPollId64)
	voterPoll := db.VoterPoll{PollID: pollId, VoteDate: time.Now()}
	c.ShouldBindJSON(&voterPoll)    // try to load json, it is fine if json is not provided
	if pollId != voterPoll.PollID { // load success, but parameter doesn't match body
		log.Println("URL parameter (pollid) does not match Request Body (voterPoll.PollID)")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// load success, matched
	if err := va.db.AddVoterPoll(uint(voterId64), voterPoll); err != nil {
		log.Println("Error adding voter poll: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, voterPoll)
}

func (va *VoteAPI) UpdateVoter(c *gin.Context) {
	var voter db.Voter
	if err := c.ShouldBindJSON(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := va.db.UpdateVoter(voter); err != nil {
		log.Println("Error updating voter: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, voter)
}

func (va *VoteAPI) DeleteVoter(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := va.db.DeleteVoter(uint(voterId64)); err != nil {
		log.Println("Error deleting voter: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func (va *VoteAPI) DeleteVoterPoll(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	voterPollIdS := c.Param("pollid")
	voterPollId64, err := strconv.ParseUint(voterPollIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := va.db.DeleteVoterPoll(uint(voterId64), uint(voterPollId64)); err != nil {
		log.Println("Error deleting voter poll: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

/*   SPECIAL HANDLERS FOR DEMONSTRATION - CRASH SIMULATION AND HEALTH CHECK */

// implementation for GET /crash
// This simulates a crash to show some of the benefits of the
// gin framework
func (td *VoteAPI) CrashSim(c *gin.Context) {
	//panic() is go's version of throwing an exception
	panic("Simulating an unexpected crash")
}

// implementation of GET /health. It is a good practice to build in a
// health check for your API.  Below the results are just hard coded
// but in a real API you can provide detailed information about the
// health of your API with a Health Check
func (td *VoteAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"status":             "ok",
			"version":            "1.0.0",
			"uptime":             100,
			"users_processed":    1000,
			"errors_encountered": 10,
		})
}
