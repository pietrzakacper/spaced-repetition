package main

import (
	"controller"
	"interactor"
)

func main() {
	var i interactor.Interactor = interactor.HttpInteractor{}

	var c controller.Controller = controller.FlashcardsController{}

	i.Start(c)
}
