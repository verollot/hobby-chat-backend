package main

import (
	"os"

	"github.com/frisk038/livechat/app/handlers"
	"github.com/frisk038/livechat/app/handlers/connexions"
	"github.com/frisk038/livechat/business"
	"github.com/frisk038/livechat/infra/repo"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func initRoutes(hp handlers.HandlerProfile, hc handlers.HandlerChat) {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/user", hp.PostUsers)
	r.POST("/user/:user_id/hobby/:hobby", hp.PostUsersHobbies)
	r.GET("/user/:user_id/hobbies", hp.GetUsersHobbies)
	r.DELETE("/user/:user_id/hobby/:hobby_id", hp.DelUsersHobbies)

	r.GET("/ws/:user_id", hc.RegisterClientSocket)

	r.Run(":" + port)
}

func main() {
	repo, err := repo.NewRepo()
	if err != nil {
		log.Error(err)
	}

	connxs := connexions.NewConnexionsMap()
	bp := business.NewBusinessProfile(repo)
	hp := handlers.NewHandlerProfile(&bp)
	hc := handlers.NewHandlerChat(connxs)

	initRoutes(hp, hc)
}
