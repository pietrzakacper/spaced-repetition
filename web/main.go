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

	if os.Getenv("TT_TOKEN") == "" {
		tt.Config.Disable()
	}

	tt.Log("start app", "hello")
	log.Println("Running interactor.Start()...")
	i.Start()
}
