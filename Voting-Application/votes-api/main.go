package main

import (
	"flag"
	"fmt"
	"os"
	"votes-api/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 80, "Default Port")

	flag.Parse()
}

func main() {
	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	apiHandler, err := api.NewVotesAPI()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r.GET("/votes", apiHandler.ListAllVotes)
	r.DELETE("/votes", apiHandler.DeleteAllVotes)

	r.GET("/votes/:id", apiHandler.GetVote)
	r.POST("/votes/:id", apiHandler.AddVote)
	r.PUT("/votes/:id", apiHandler.UpdateVote)
	r.PUT("/votes/", apiHandler.UpdateVote)
	r.DELETE("/votes/:id", apiHandler.DeleteVote)

	r.GET("/votes/health", apiHandler.HealthCheck)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
