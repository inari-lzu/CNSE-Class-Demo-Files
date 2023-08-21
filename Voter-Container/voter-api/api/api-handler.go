package api

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"voter-api/db"

	"github.com/gin-gonic/gin"
)

type VoteAPI struct {
	voterList   *db.VoterList
	bootTime    time.Time
	successes   int
	badRequests int
}

func NewVoteAPI() (*VoteAPI, error) {
	dbHandler, err := db.NewVoterList()
	if err != nil {
		return nil, err
	}

	return &VoteAPI{voterList: dbHandler, bootTime: time.Now(), successes: 0, badRequests: 0}, nil
}

func (va *VoteAPI) ListAllVoters(c *gin.Context) {

	vl, err := va.voterList.GetAllVoters()
	if err != nil {
		log.Println("Error Getting All Voters: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if vl == nil {
		vl = make([]db.Voter, 0)
	}

	c.JSON(http.StatusOK, vl)
	va.successes++
}

func (va *VoteAPI) DeleteAllVoters(c *gin.Context) {
	err := va.voterList.DeleteAllVoters()
	if err != nil {
		log.Println("Error deleting All Voters: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
	va.successes++
}

func (va *VoteAPI) GetVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	voter, err := va.voterList.GetVoter(uint(id64))
	if err != nil {
		log.Println("voter not found: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voter)
	va.successes++
}

func (va *VoteAPI) AddVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var voter db.Voter
	if err := c.ShouldBindJSON(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if uint(id64) != voter.VoterID {
		log.Println("URL parameter (id) does not match Request Body (voter.VoterID)")
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := va.voterList.AddVoter(voter); err != nil {
		log.Println("Error adding voter: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, voter)
	va.successes++
}

func (va *VoteAPI) UpdateVoter(c *gin.Context) {
	var v db.Voter
	if err := c.ShouldBindJSON(&v); err != nil {
		log.Println("Error binding JSON: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	voter, err := va.voterList.UpdateVoter(v)
	if err != nil {
		log.Println("Error updating voter: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, voter)
	va.successes++
}

func (va *VoteAPI) GetVoterHistory(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	voteHistory, err := va.voterList.GetVoteHistory(uint(id64))
	if err != nil {
		log.Println("Get voteHistory fail: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voteHistory)
	va.successes++
}

func (va *VoteAPI) DeleteVoterHistory(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = va.voterList.DeleteVoteHistory(uint(id64))
	if err != nil {
		log.Println("Error deleting voteHistory: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Status(http.StatusOK)
	va.successes++
}

func (va *VoteAPI) GetVoterPoll(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	voterPollIdS := c.Param("pollid")
	voterPollId64, err := strconv.ParseUint(voterPollIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	voterPoll, err := va.voterList.GetVoterPoll(uint(voterId64), uint(voterPollId64))
	if err != nil {
		log.Println("voter poll not found: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, voterPoll)
	va.successes++
}

func (va *VoteAPI) AddVoterPoll(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var v db.Voter
	if err := c.ShouldBindJSON(&v); err != nil {
		log.Println("Error binding JSON: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if len(v.VoteHistory) != 1 {
		log.Println("len(v.VoteHistory) != 1")
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	vp := v.VoteHistory[0]
	if err := va.voterList.AddVoterPoll(uint(voterId64), vp); err != nil {
		log.Println("Error adding voter poll: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, vp)
	va.successes++
}

func (va *VoteAPI) UpdateVoterPoll(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var v db.Voter
	if err := c.ShouldBindJSON(&v); err != nil {
		log.Println("Error binding JSON: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if len(v.VoteHistory) != 1 {
		log.Println("len(v.VoteHistory) != 1")
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	vp := v.VoteHistory[0]
	if err := va.voterList.UpdateVoterPoll(uint(voterId64), vp); err != nil {
		log.Println("Error updating voter poll: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, vp)
	va.successes++
}

func (va *VoteAPI) DeleteVoter(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := va.voterList.DeleteVoter(uint(voterId64)); err != nil {
		log.Println("Error deleting voter: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	va.successes++
}

func (va *VoteAPI) DeleteVoterPoll(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	voterPollIdS := c.Param("pollid")
	voterPollId64, err := strconv.ParseUint(voterPollIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := va.voterList.DeleteVoterPoll(uint(voterId64), uint(voterPollId64)); err != nil {
		log.Println("Error deleting voter poll: ", err)
		va.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	va.successes++
}

func (va *VoteAPI) CrashSim(c *gin.Context) {
	va.badRequests++
	panic("Simulating an unexpected crash")
}

func (va *VoteAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"status":                      http.StatusOK,
			"gin_version":                 gin.Version,
			"api_uptime":                  time.Since(va.bootTime).String(),
			"total_api_calls":             va.successes + va.badRequests,
			"total_api_calls_succeed":     va.successes,
			"total_api_calls_with_errors": va.badRequests,
		},
	)
	va.successes++
}
