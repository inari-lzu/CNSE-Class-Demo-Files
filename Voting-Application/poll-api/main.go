package main

import (
	"flag"
	"fmt"
	"os"
	"poll-api/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1081, "Default Port")

	flag.Parse()
}

func main() {
	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	apiHandler, err := api.NewPollAPI()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r.GET("/polls", apiHandler.ListAllPolls)
	r.DELETE("/polls", apiHandler.DeleteAllPolls)

	r.GET("/polls/:id", apiHandler.GetPoll)
	r.POST("/polls/:id", apiHandler.AddPoll)
	r.PUT("/polls/:id", apiHandler.UpdatePoll)
	r.PUT("/polls/", apiHandler.UpdatePoll)
	r.DELETE("/polls/:id", apiHandler.DeletePoll)

	r.GET("/polls/:id/options", apiHandler.GetAllOptions)
	r.DELETE("/polls/:id/options", apiHandler.DeleteAllOptions)

	r.GET("/polls/:id/options/:optionid", apiHandler.GetOption)
	r.POST("/polls/:id/options/:optionid", apiHandler.AddOption)
	r.POST("/polls/:id/options/", apiHandler.AddOption)
	r.PUT("/polls/:id/options/:optionid", apiHandler.UpdateOption)
	r.PUT("/polls/:id/options/", apiHandler.UpdateOption)
	r.DELETE("/polls/:id/options/:optionid", apiHandler.DeleteOption)

	r.GET("/polls/health", apiHandler.HealthCheck)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
