package main

import (
	"fmt"
	"persistance"
	"user"
	"web/interactor"
)

func main() {
	fmt.Println("Running main...")
	var userSessionFactory = &user.UserSessionFactory{Persistance: &persistance.PostgresPersistance{}}
	var i interactor.Interactor = interactor.CreateHttpInteractor(userSessionFactory)

	fmt.Println("Running interactor.Start()...")
	i.Start()
}
