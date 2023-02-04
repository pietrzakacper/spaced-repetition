package main

import (
	"controller"
	"persistance"
	"web/interactor"
	"web/view"
)

func main() {
	var httpView = &view.HttpView{}
	var i interactor.Interactor = interactor.CreateHttpInteractor(httpView)

	var v controller.View = httpView
	var p controller.Persistance = &persistance.BadgerPersistance{}
	var flashcardController = controller.CreateFlashcardsController(v, p)
	var c controller.Controller = flashcardController

	i.Start(c)
}
