package api

import (
	"log"
	"net/http"
	"poll-api/poll"
	"strconv"
	"time"

	"db"

	"github.com/gin-gonic/gin"
)

type PollAPI struct {
	polls       *db.Handler[poll.Poll]
	bootTime    time.Time
	successes   int
	badRequests int
}

func NewPollAPI() (*PollAPI, error) {
	dbHandler, err := db.NewHandler[poll.Poll]("poll")
	if err != nil {
		return nil, err
	}

	return &PollAPI{polls: dbHandler, bootTime: time.Now(), successes: 0, badRequests: 0}, nil
}

func (api *PollAPI) ListAllPolls(c *gin.Context) {
	pollList, err := api.polls.All()
	if err != nil {
		log.Println("Error Getting All Polls: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if pollList == nil {
		pollList = make([]poll.Poll, 0)
	}
	urlList := make([]string, 0)
	for _, p := range pollList {
		urlList = append(urlList, p.ToJson().Poll)
	}
	c.JSON(http.StatusOK, urlList)
	api.successes++
}

func (api *PollAPI) DeleteAllPolls(c *gin.Context) {
	err := api.polls.Clear()
	if err != nil {
		log.Println("Error deleting All Polls: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
	api.successes++
}

func (api *PollAPI) GetPoll(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	p, err := api.polls.Get(uint(id64))
	if err != nil {
		log.Println("poll not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, p.ToJson())
	api.successes++
}

func (api *PollAPI) AddPoll(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	p := poll.NewPoll()
	if err := c.ShouldBindJSON(&p); err != nil {
		log.Println("Error binding JSON: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if uint(id64) != p.PollID {
		log.Println("URL parameter (id) does not match Request Body (poll.PollID)")
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = api.polls.Add(p)
	if err != nil {
		log.Println("Error adding poll: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, p.ToJson())
	api.successes++
}

func (api *PollAPI) UpdatePoll(c *gin.Context) {
	var p poll.Poll
	if err := c.ShouldBindJSON(&p); err != nil {
		log.Println("Error binding JSON: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	p, err := api.polls.Update(p, poll.Poll.Update)
	if err != nil {
		log.Println("Error updating PollInfo: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, p.ToJson())
	api.successes++
}

func (api *PollAPI) DeletePoll(c *gin.Context) {
	pollIdS := c.Param("id")
	pollId64, err := strconv.ParseUint(pollIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := api.polls.Delete(uint(pollId64)); err != nil {
		log.Println("Error deleting poll: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	api.successes++
}

func (api *PollAPI) GetAllOptions(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	p, err := api.polls.Get(uint(id64))
	if err != nil {
		log.Println("poll not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, p.ToJson().Options)
	api.successes++
}

func (api *PollAPI) DeleteAllOptions(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	p, err := api.polls.Get(uint(id64))
	if err != nil {
		log.Println("poll not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	_, err = api.polls.Update(p.DeleteAllOptions(), poll.Poll.UpdateOptions)
	if err != nil {
		log.Println("Error updating PollOptions: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	api.successes++
}

func (api *PollAPI) GetOption(c *gin.Context) {
	pollIdS := c.Param("id")
	pollId64, err := strconv.ParseUint(pollIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	optionIdS := c.Param("optionid")
	optionId64, err := strconv.ParseUint(optionIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	p, err := api.polls.Get(uint(pollId64))
	if err != nil {
		log.Println("poll not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	option, err := p.GetOption(uint(optionId64))
	if err != nil {
		log.Println("poll option not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, p.GetOptionJson(option))
	api.successes++
}

func (api *PollAPI) AddOption(c *gin.Context) {
	pollIdS := c.Param("id")
	pollId64, err := strconv.ParseUint(pollIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var p poll.Poll
	if err := c.ShouldBindJSON(&p); err != nil {
		log.Println("Error binding JSON: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if len(p.PollOptions) != 1 {
		log.Println("len(p.PollOptions) != 1")
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	po := p.PollOptions[0]

	p, err = api.polls.Get(uint(pollId64))
	if err != nil {
		log.Println("poll not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	updatedPoll, err := p.AddPoll(po)
	if err != nil {
		log.Println("Error adding poll option: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = api.polls.Update(updatedPoll, poll.Poll.UpdateOptions)
	if err != nil {
		log.Println("Error updating PollOptions: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, p.GetOptionJson(po))
	api.successes++
}

func (api *PollAPI) UpdateOption(c *gin.Context) {
	pollIdS := c.Param("id")
	pollId64, err := strconv.ParseUint(pollIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var p poll.Poll
	if err := c.ShouldBindJSON(&p); err != nil {
		log.Println("Error binding JSON: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if len(p.PollOptions) != 1 {
		log.Println("len(p.PollOptions) != 1")
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	po := p.PollOptions[0]

	p, err = api.polls.Get(uint(pollId64))
	if err != nil {
		log.Println("poll not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	updatedPoll, err := p.UpdateOption(po)
	if err != nil {
		log.Println("Error updating poll option: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = api.polls.Update(updatedPoll, poll.Poll.UpdateOptions)
	if err != nil {
		log.Println("Error updating pollOptions: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, p.GetOptionJson(po))
	api.successes++
}

func (api *PollAPI) DeleteOption(c *gin.Context) {
	pollIdS := c.Param("id")
	pollId64, err := strconv.ParseUint(pollIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	optionIdS := c.Param("optionid")
	optionId64, err := strconv.ParseUint(optionIdS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	p, err := api.polls.Get(uint(pollId64))
	if err != nil {
		log.Println("poll not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	updatedPoll, err := p.DeleteOption(uint(optionId64))
	if err != nil {
		log.Println("poll option not found: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = api.polls.Update(updatedPoll, poll.Poll.UpdateOptions)
	if err != nil {
		log.Println("Error updating pollOptions: ", err)
		api.badRequests++
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
	api.successes++
}

func (api *PollAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"api_name":                    "poll-api",
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
