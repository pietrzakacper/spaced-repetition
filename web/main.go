package main

import (
	"persistance"
	"user"
	"web/interactor"
	"web/view"
)

// @TODO make the code threadsafe
func main() {
	var httpView = &view.HttpView{}
	var userSessionFactory = &user.UserSessionFactory{View: httpView, Persistance: &persistance.BadgerPersistance{}}
	var i interactor.Interactor = interactor.CreateHttpInteractor(httpView, userSessionFactory)

	i.Start()
}
