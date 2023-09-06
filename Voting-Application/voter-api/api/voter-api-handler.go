package api

import (
	"db"
	"log"
	"net/http"
	"strconv"
	"time"
	"voter-api/voter"

	"github.com/gin-gonic/gin"
)

type VoterAPI struct {
	voters      *db.Handler[voter.Voter]
	bootTime    time.Time
	successes   int
	badRequests int
}

func NewVoterAPI() (*VoterAPI, error) {
	dbHandler, err := db.NewHandler[voter.Voter]("voter")
	if err != nil {
		return nil, err
	}

	return &VoterAPI{voters: dbHandler, bootTime: time.Now(), successes: 0, badRequests: 0}, nil
}

func (api *VoterAPI) ListAllVoters(c *gin.Context) {
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

	urlList := make([]string, 0)
	for _, vr := range voterList {
		urlList = append(urlList, vr.ToJson().Voter)
	}

	c.JSON(http.StatusOK, urlList)
	api.successes++
}

func (api *VoterAPI) DeleteAllVoters(c *gin.Context) {
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

func (api *VoterAPI) GetVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	vr, err := api.voters.Get(uint(id64))
	if err != nil {
		log.Println("voter not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, vr.ToJson())
	api.successes++
}

func (api *VoterAPI) AddVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	vr := voter.NewVoter()
	if err := c.ShouldBindJSON(&vr); err != nil {
		log.Println("Error binding JSON: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if uint(id64) != vr.VoterID {
		log.Println("URL parameter (id) does not match Request Body (voter.VoterID)")
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = api.voters.Add(vr.DeleteHistory())
	if err != nil {
		log.Println("Error adding voter: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, vr.ToJson())
	api.successes++
}

func (api *VoterAPI) DeleteVoter(c *gin.Context) {
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

func (api *VoterAPI) UpdateVoter(c *gin.Context) {
	var vr voter.Voter
	if err := c.ShouldBindJSON(&vr); err != nil {
		log.Println("Error binding JSON: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	vr, err := api.voters.Update(vr, voter.Voter.UpdateVoterInfo)
	if err != nil {
		log.Println("Error updating VoterInfo: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, vr.ToJson())
	api.successes++
}

func (api *VoterAPI) GetVoterHistory(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	vr, err := api.voters.Get(uint(id64))
	if err != nil {
		log.Println("voter not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, vr.ToJson().History)
	api.successes++
}

func (api *VoterAPI) DeleteVoterHistory(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	vr, err := api.voters.Get(uint(id64))
	if err != nil {
		log.Println("voter not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	_, err = api.voters.Update(vr.DeleteHistory(), voter.Voter.UpdateVoteHistory)
	if err != nil {
		log.Println("Error updating voteHistory: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	api.successes++
}

func (api *VoterAPI) GetVoterPoll(c *gin.Context) {
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

	vr, err := api.voters.Get(uint(voterId64))
	if err != nil {
		log.Println("voter not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	vp, err := vr.GetPoll(uint(voterPollId64))
	if err != nil {
		log.Println("voter poll not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, vr.GetPollJson(vp))
	api.successes++
}

func (api *VoterAPI) AddVoterPoll(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var vr voter.Voter
	if err := c.ShouldBindJSON(&vr); err != nil {
		log.Println("Error binding JSON: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if len(vr.VoteHistory) != 1 {
		log.Println("len(vr.VoteHistory) != 1")
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	vp := vr.VoteHistory[0]

	vr, err = api.voters.Get(uint(voterId64))
	if err != nil {
		log.Println("voter not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	updatedVoter, err := vr.AddPoll(vp)
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

	c.JSON(http.StatusOK, vr.GetPollJson(vp))
	api.successes++
}

func (api *VoterAPI) UpdateVoterPoll(c *gin.Context) {
	voterIdS := c.Param("id")
	voterId64, err := strconv.ParseUint(voterIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var vr voter.Voter
	if err := c.ShouldBindJSON(&vr); err != nil {
		log.Println("Error binding JSON: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if len(vr.VoteHistory) != 1 {
		log.Println("len(vr.VoteHistory) != 1")
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	vp := vr.VoteHistory[0]

	vr, err = api.voters.Get(uint(voterId64))
	if err != nil {
		log.Println("voter not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	updatedVoter, err := vr.UpdatePoll(vp)
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

	c.JSON(http.StatusOK, vr.GetPollJson(vp))
	api.successes++
}

func (api *VoterAPI) DeleteVoterPoll(c *gin.Context) {
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

	vr, err := api.voters.Get(uint(voterId64))
	if err != nil {
		log.Println("voter not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	updatedVoter, err := vr.DeletePoll(uint(voterPollId64))
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

func (api *VoterAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"api_name":                    "voter-api",
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
