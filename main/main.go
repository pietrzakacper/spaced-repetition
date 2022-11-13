package main

import (
	"controller"
	"interactor"
	"view"
)

func main() {
	var httpView = &view.HttpView{}

	var i interactor.Interactor = interactor.CreateHttpInteractor(httpView)

	var v controller.View = httpView
	var flashcardController = controller.CreateFlashcardsController(v)
	var c controller.Controller = flashcardController

	i.Start(c)
}
