package user

import "controller"

type UserSessionFactory struct {
	controller.View
	controller.Persistance
}

func (u *UserSessionFactory) Create(userId string) *controller.FlashcardsController {
	store := u.Persistance.Create("db", userId)
	// @TODO learn why we cannot return a pointer here
	return controller.CreateFlashcardsController(u.View, store)
}

// var flashcardController = controller.CreateFlashcardsController(v, p)

// var v controller.View = httpView
// var p controller.Persistance = &persistance.CSVPersistance{}

// type UserSession interface {
// 	GetController() Controller
// }

// type Persistance interface {
// 	Create(name string) StoreFactory
// }

// type StoreFactory interface {
// 	Create(userId string) Store
// }
