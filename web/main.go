package main

import (
	"persistance"
	"user"
	"web/interactor"
)

// @TODO make the code threadsafe
func main() {
	var userSessionFactory = &user.UserSessionFactory{Persistance: &persistance.CSVPersistance{}}
	var i interactor.Interactor = interactor.CreateHttpInteractor(userSessionFactory)

	i.Start()
}
