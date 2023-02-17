package main

import (
	"persistance"
	"user"
	"web/interactor"
	"web/view"
)

func main() {
	var httpView = &view.HttpView{}
	var userSessionFactory = &user.UserSessionFactory{View: httpView, Persistance: &persistance.CSVPersistance{}}
	var i interactor.Interactor = interactor.CreateHttpInteractor(httpView, userSessionFactory)

	i.Start()
}
