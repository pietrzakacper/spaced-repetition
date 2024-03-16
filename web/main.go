package main

import (
	"log"
	"os"
	"persistance"
	"user"
	"web/interactor"

	"github.com/pietrzakacper/tracethat.dev/reporters/golang/tt"
)

func main() {
	log.Println("Running main...")
	var userSessionFactory = &user.UserSessionFactory{Persistance: &persistance.PostgresPersistance{}}
	var i interactor.Interactor = interactor.CreateHttpInteractor(userSessionFactory)

	if ttToken := os.Getenv("TT_TOKEN"); ttToken != "" {
		tt.RegisterToken(ttToken)
	} else {
		tt.DisableDevtools()
	}
	if ttServerUrl := os.Getenv("TT_SERVER_URL"); ttServerUrl != "" {
		tt.SetServerUrl(ttServerUrl)
	}
	tt.Log("start app", "hello")
	log.Println("Running interactor.Start()...")
	i.Start()
}
