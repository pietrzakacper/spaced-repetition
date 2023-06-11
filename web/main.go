package main

import (
	"log"
	"persistance"
	"user"
	"web/interactor"
)

func main() {
	log.Println("Running main...")
	var userSessionFactory = &user.UserSessionFactory{Persistance: &persistance.PostgresPersistance{}}
	var i interactor.Interactor = interactor.CreateHttpInteractor(userSessionFactory)

	log.Println("Running interactor.Start()...")
	i.Start()
}
