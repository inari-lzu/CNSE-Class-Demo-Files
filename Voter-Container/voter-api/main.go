package main

import (
	"flag"
	"fmt"
	"os"
	"voter-api/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

func main() {
	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	apiHandler, err := api.NewVoteAPI()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r.GET("/voters", apiHandler.ListAllVoters)
	r.DELETE("/voters", apiHandler.DeleteAllVoters)

	r.GET("/voters/:id", apiHandler.GetVoter)
	r.POST("/voters/:id", apiHandler.AddVoter)
	r.PUT("/voters/:id", apiHandler.UpdateVoter)
	r.DELETE("/voters/:id", apiHandler.DeleteVoter)

	r.GET("/voters/:id/polls", apiHandler.GetVoterHistory)
	r.DELETE("/voters/:id/polls", apiHandler.DeleteVoterHistory)

	r.GET("/voters/:id/polls/:pollid", apiHandler.GetVoterPoll)
	r.POST("/voters/:id/polls/", apiHandler.AddVoterPoll)
	r.PUT("/voters/:id/polls/", apiHandler.UpdateVoterPoll)
	r.DELETE("/voters/:id/polls/:pollid", apiHandler.DeleteVoterPoll)

	r.GET("/voters/health", apiHandler.HealthCheck)

	v2 := r.Group("/v2")
	v2.GET("/crash", apiHandler.CrashSim)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
